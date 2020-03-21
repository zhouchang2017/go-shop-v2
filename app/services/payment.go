package services

import (
	"context"
	"errors"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/wechat"
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
	if order.User.Id != userInfo.ID.String() {
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
		Title:          fmt.Sprint("订单号%s", order.OrderNo),
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
		Body:           fmt.Sprint("订单号%s", order.OrderNo),
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
func (srv *PaymentOption) Notify(ctx context.Context) {

}
