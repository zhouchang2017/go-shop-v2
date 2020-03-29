package mp_subscribe

import (
	"fmt"
	"github.com/medivhzhan/weapp/v2"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/utils"
)

// 用户订单付款通知
type OrderPaidNotify struct {
	order *models.Order
}

func NewOrderPaidNotify(order *models.Order) *OrderPaidNotify {
	return &OrderPaidNotify{order: order}
}

func (o OrderPaidNotify) To() string {
	return o.order.User.WechatMiniId
}

func (o OrderPaidNotify) TemplateID() string {
	return "tr7pmPfvX5NxO_uFtSDbL1vi2Pk1U9_H4oCOEC4MlQw"
}

func (o OrderPaidNotify) Page() string {
	return fmt.Sprintf("pages/home/order/detail?id=%s", o.order.GetID())
}

func (o OrderPaidNotify) Data() weapp.SubscribeMessageData {
	return weapp.SubscribeMessageData{
		"character_string6": {Value: o.order.OrderNo},                               // 订单编号
		"thing1":            {Value: o.order.GoodsName(20)},                         // 商品名称
		"amount3":           {Value: o.order.GetActualAmount()},                     // 支付金额
		"date5":             {Value: utils.TimeJsonOut(*o.order.Payment.PaymentAt)}, // 下单时间
		"thing10":           {Value: "您以成功下单，我们会尽快为您安排发货"},                          // 温馨提示
	}
}
