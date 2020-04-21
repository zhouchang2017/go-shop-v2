package mp_subscribe

import (
	"fmt"
	"github.com/medivhzhan/weapp/v2"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/utils"
)

// 退款变动提醒
type OrderRefundChangeNotify struct {
	refund *models.Refund
}

func NewOrderRefundChangeNotify(refund *models.Refund) *OrderRefundChangeNotify {
	return &OrderRefundChangeNotify{refund}
}

func (o OrderRefundChangeNotify) To() string {
	return o.refund.OpenId
}

func (o OrderRefundChangeNotify) TemplateID() string {
	return "c_r-6vdE4Pg0UYICu_HcRRhPhY1WzUD_Fd6oNx2KMes"
}

func (o OrderRefundChangeNotify) Page() string {
	return fmt.Sprintf("pages/home/order/detail?id=%s", o.refund.OrderId)
}

func (o OrderRefundChangeNotify) Data() weapp.SubscribeMessageData {
	return weapp.SubscribeMessageData{
		"character_string4": {Value: o.refund.OrderNo},
		"thing5":            {Value: o.refund.GoodsName(20)},
		"thing1":            {Value: utils.SubString(o.refund.RefundDesc, 0, 20)},
		"amount2":           {Value: o.refund.GetActualAmount()},
		"thing6":            {Value: o.refund.StatusText()},
	}
}
