package mp_subscribe

import (
	"fmt"
	"github.com/medivhzhan/weapp/v2"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/utils"
)

// 退款变动提醒
type OrderRefundChangeNotify struct {
	order    *models.Order
	refundId string
}

func NewOrderRefundChangeNotify(order *models.Order, refundId string) *OrderRefundChangeNotify {
	return &OrderRefundChangeNotify{order: order, refundId: refundId}
}

func (o OrderRefundChangeNotify) To() string {
	return o.order.User.WechatMiniId
}

func (o OrderRefundChangeNotify) TemplateID() string {
	return "c_r-6vdE4Pg0UYICu_HcRRhPhY1WzUD_Fd6oNx2KMes"
}

func (o OrderRefundChangeNotify) Page() string {
	return fmt.Sprintf("pages/home/order/detail?id=%s", o.order.GetID())
}

func (o OrderRefundChangeNotify) Data() weapp.SubscribeMessageData {
	refund, _ := o.order.FindRefund(o.refundId)
	return weapp.SubscribeMessageData{
		"character_string4": {Value: o.order.OrderNo},
		"thing5":            {Value: refund.GoodsName(20, o.order)},
		"thing1":            {Value: utils.SubString(refund.RefundDesc, 0, 20)},
		"amount2":           {Value: refund.GetActualAmount()},
		"thing6":            {Value: refund.StatusText()},
	}
}
