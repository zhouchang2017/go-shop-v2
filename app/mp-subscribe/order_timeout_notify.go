package mp_subscribe

import (
	"fmt"
	"github.com/medivhzhan/weapp/v2"
	"go-shop-v2/app/models"
)

// 订单关闭提醒
type OrderTimeoutNotify struct {
	order *models.Order
}

func NewOrderClosedNotify(order *models.Order) *OrderTimeoutNotify {
	return &OrderTimeoutNotify{order: order}
}

func (o OrderTimeoutNotify) To() string {
	return o.order.User.WechatMiniId
}

func (o OrderTimeoutNotify) TemplateID() string {
	return "gKqinsdq2MiIPjF_bnIhotZ0DlrRcnHb6Ugax25dXJU"
}

func (o OrderTimeoutNotify) Page() string {
	return fmt.Sprintf("pages/home/order/detail?id=%s", o.order.GetID())
}

func (o OrderTimeoutNotify) Data() weapp.SubscribeMessageData {
	return weapp.SubscribeMessageData{
		"character_string3": {Value: o.order.OrderNo},
		"thing4":            {Value: o.order.GoodsName(20)},
		"amount2":           {Value: o.order.GetActualAmount()},
		"phrase1":           {Value: o.order.StatusText()},
		"thing5":            {Value: *o.order.CloseReason},
	}
}
