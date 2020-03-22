package services

import (
	"context"
	"errors"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/wechat"
	gopay_wechat "github.com/iGoogle-ink/gopay/wechat"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type PaymentService struct {
	paymentRep *repositories.PaymentRep
	orderResp  *repositories.OrderRep
}

func NewPaymentService(paymentRep *repositories.PaymentRep, orderResp *repositories.OrderRep) *PaymentService {
	return &PaymentService{paymentRep: paymentRep, orderResp: orderResp}
}

type PaymentOption struct {
	OrderId        string `json:"order_id" form:"order_id" binding:"required"`
	OrderNo        string `json:"order_no" form:"order_no" binding:"required"`
	SpbillCreateIp string // 用户ip
	Platform       string `json:"platform"` // 支付平台:微信/支付宝
}

func (p PaymentOption) GetPlatform() string {
	if p.Platform == "" {
		return "wechat"
	}
	return p.Platform
}

func (paymentOpt *PaymentOption) IsValid() error {
	if paymentOpt.OrderId == "" {
		return errors.New("OrderId is empty")
	}
	if paymentOpt.OrderNo == "" {
		return errors.New("OrderNo is empty")
	}
	return nil
}

// 下单
func (srv *PaymentService) Payment(ctx context.Context, userInfo *models.User, opt *PaymentOption) (*wechat.WechatMiniPayConfig, error) {
	if err := opt.IsValid(); err != nil {
		return nil, err
	}
	// get order information
	orderRes := <-srv.orderResp.FindById(ctx, opt.OrderId)
	if orderRes.Error != nil {
		return nil, orderRes.Error
	}
	order := orderRes.Result.(*models.Order)
	if order.OrderNo != opt.OrderNo {
		return nil, errors.New("invalid OrderNo with different")
	}
	if order.User.Id != userInfo.GetID() {
		return nil, errors.New("invalid permission caused of different user")
	}
	if order.Status != models.OrderStatusPrePay {
		return nil, errors.New("invalid order status not prepay")
	}

	// unify order
	option := srv.setUnifiedOrderOption(order, userInfo.WechatMiniId, opt.SpbillCreateIp)
	wxRsp, wcErr := wechat.Pay.UnifiedOrder(option)

	if wcErr != nil {
		return nil, wcErr
	}

	// update payment information in db
	srv.paymentRep.Store(ctx, &models.Payment{
		OrderNo:        order.OrderNo,
		Platform:       opt.GetPlatform(),
		Title:          fmt.Sprintf("订单号%s", order.OrderNo),
		Amount:         order.ActualAmount,
		ExtendedUserId: userInfo.GetID(),
		PrePaymentNo:   wxRsp.PrepayId,
		PaymentNo:      "",
		PaymentAt:      nil,
	})
	return wxRsp.GetWechatMiniPayConfig(), nil
}

func (srv *PaymentService) setUnifiedOrderOption(order *models.Order, openId string, spbillCreateIp string) (opt *wechat.PayUnifiedOrderOption) {

	opt = &wechat.PayUnifiedOrderOption{
		Body:           fmt.Sprintf("订单号%s", order.OrderNo),
		Detail:         nil,
		OutTradeNo:     order.OrderNo,
		TotalFee:       order.ActualAmount,
		SpbillCreateIp: spbillCreateIp,
		OpenId:         openId,
	}
	if len(order.OrderItems) > 1 {
		detail := &wechat.PayOrderDetail{
			CostPrice:   order.ActualAmount,
			GoodsDetail: nil,
		}
		for _, item := range order.OrderItems {
			detail.GoodsDetail = append(detail.GoodsDetail, &wechat.PayOrderGoodsDetail{
				GoodsId:   item.Item.Id,
				GoodsName: item.Item.Product.Name,
				Quantity:  item.Count,
				Price:     item.Amount,
			})
		}
		opt.Detail = detail
	}
	return opt
}

// 支付回调
func (srv *PaymentService) PayNotify(ctx context.Context, req *http.Request) error {
	// parse
	notifyReq, notifyErr := gopay_wechat.ParseNotifyResult(req)
	if notifyErr != nil {
		return fmt.Errorf("parse notify result failed with %s", notifyErr)
	}
	// check sign
	verifyOk, verifyErr := gopay_wechat.VerifySign(wechat.Pay.ApiKey, gopay_wechat.SignType_MD5, notifyReq)
	if verifyErr != nil {
		return fmt.Errorf("verify sign occur error %s", verifyErr)
	}
	if !verifyOk {
		return errors.New("verify sign failed")
	}
	// deal with
	orderNo := notifyReq.OutTradeNo
	paymentNo := notifyReq.TransactionId
	// get order information
	orderRes := <-srv.orderResp.FindOne(ctx, map[string]interface{}{
		"order_no": orderNo,
	})
	if orderRes.Error != nil {
		return orderRes.Error
	}
	order := orderRes.Result.(*models.Order)
	order.Payment.PaymentNo = paymentNo
	// check status (这里更新了中间态，但是目前不保存先)
	switch order.Status {
	case models.OrderStatusPrePay:
		order.Status = models.OrderStatusPaid
	case models.OrderStatusPaid:
		//continue
	default:
		//log.Println("already final status")
		return nil
	}
	// check result
	if notifyReq.ReturnCode == "SUCCESS" {
		// if failed
		if notifyReq.ResultCode != "SUCCESS" {
			// update payment
			updateRes := <-srv.orderResp.Update(ctx, order.GetID(), bson.M{
				"$set": bson.M{
					"status": models.OrderStatusFailed,
					"payment.payment_no": paymentNo,
				},
			})
			if updateRes.Error != nil {
				return fmt.Errorf("update order-%s to failed status failed %s", orderNo, updateRes.Error)
			}
			// return
			//log.Println("update order xxx to failed status success")
			return nil
		}
		// deal with if success
		// valid money
		if uint64(notifyReq.TotalFee) != order.ActualAmount {
			return fmt.Errorf("notify total fee %d is different with order actucal amount %d", notifyReq.TotalFee, order.ActualAmount)
		}
		// update order
		updateRes := <-srv.orderResp.Update(ctx, order.GetID(), bson.M{
			"$set": bson.M{
				"status": models.OrderStatusPreSend,
				"payment.payment_no": paymentNo,
			},
		})
		if updateRes.Error != nil {
			return fmt.Errorf("update order-%s to success status failed %s", orderNo, updateRes.Error)
		}
		// todo: add other case to do and remember use transaction if involving other table in db

		// return
		//log.Println("update order xxx to success status success")
		return nil
	} else {
		return errors.New("return code not SUCCESS")
	}
	// return
	return nil
}
