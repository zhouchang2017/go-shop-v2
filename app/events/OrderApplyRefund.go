package events

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"time"
)

// 订单申请退款
type OrderApplyRefund struct {
	Order    *models.Order `json:"order"`
	RefundId string        `json:"refund_id"`
}

func (o OrderApplyRefund) Delay() time.Duration {
	return 0
}

func (o OrderApplyRefund) Body() []byte {
	marshal, err := json.Marshal(o)
	if err != nil {
		log.Errorf("OrderApplyRefund json marshal error:%s", err)
	}
	return marshal
}

func NewOrderApplyRefundEvent(order *models.Order, refundId string) *OrderApplyRefund {
	return &OrderApplyRefund{Order: order, RefundId: refundId}
}
