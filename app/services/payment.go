package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/wechat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
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
func (srv *PaymentService) PayNotify(ctx context.Context, req *http.Request) (orderNumber string, err error) {

	// parse
	notifyReq, err := wechat.Pay.ParseNotifyResult(req)
	if err != nil {
		return "", err
	}
	spew.Dump(notifyReq)
	// deal with
	orderNo := notifyReq.OutTradeNo
	paymentNo := notifyReq.TransactionId

	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return "", err
	}
	if err = session.StartTransaction(); err != nil {
		return "", err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// get order information
		order, err := srv.orderResp.FindByOrderNo(sessionContext, orderNo)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}
		// 从payments中查询order对应的payment
		payment, err := srv.paymentRep.FindByOrderId(sessionContext, orderNo)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}
		if payment.PaymentAt != nil && order.Status == models.OrderStatusPaid {
			// 已经标记为支付！！
			return nil
		}
		payment.PaymentNo = paymentNo

		// check result
		if notifyReq.ReturnCode == "SUCCESS" {
			// if failed
			if notifyReq.ResultCode != "SUCCESS" {
				// 通知结果支付失败
				// update payment
				order.Payment = payment.ToAssociated()
				order.StatusToFailed()

				updateRes := <-srv.orderResp.Update(sessionContext, order.GetID(), bson.M{
					"$set": bson.M{
						"status":  models.OrderStatusFailed,
						"payment": payment.ToAssociated(),
					},
				})
				if updateRes.Error != nil {
					session.AbortTransaction(sessionContext)
					return fmt.Errorf("update order-%s to failed status failed %s", orderNo, updateRes.Error)
				}
				// return
				//log.Println("update order xxx to failed status success")
				return nil
			}
			// deal with if success
			// valid money
			if uint64(notifyReq.TotalFee) != order.ActualAmount {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("notify total fee %d is different with order actucal amount %d", notifyReq.TotalFee, order.ActualAmount)
			}
			// 支付成功
			// update payment
			payment.SetPaymentAt(time.Now())
			saved := <-srv.paymentRep.Save(sessionContext, payment)
			if saved.Error != nil {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("update order[%s] payment to paided failed %s", orderNo, saved.Error)
			}
			updatePayment := saved.Result.(*models.Payment)
			// update order
			updateRes := <-srv.orderResp.Update(sessionContext, order.GetID(), bson.M{
				"$set": bson.M{
					"status":  models.OrderStatusPreSend,
					"payment": updatePayment.ToAssociated(),
				},
			})
			if updateRes.Error != nil {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("update order-%s to success status failed %s", orderNo, updateRes.Error)
			}
			// return
			//log.Println("update order xxx to success status success")
			session.CommitTransaction(sessionContext)

			// 新标记订单为支付状态，返回订单号
			orderNumber = orderNo
			return nil
		} else {
			session.AbortTransaction(sessionContext)
			return errors.New("return code not SUCCESS")
		}
	})

	session.EndSession(ctx)
	// return
	return orderNumber, err
}

// 当然收款总金额
func (srv *PaymentService) TodayPaymentCount(ctx context.Context) (response *models.DayPaymentCount, err error) {
	return srv.paymentRep.GetRangePaymentCount(ctx, utils.TodayStart(), utils.TodayEnd())
}

// 一段时间收款聚合
func (srv *PaymentService) RangePaymentCounts(ctx context.Context, start time.Time, end time.Time) (response []*models.DayPaymentCount, err error) {
	return srv.paymentRep.GetRangePaymentCounts(ctx, start, end)
}
