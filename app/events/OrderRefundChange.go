package events

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"time"
)

type OrderRefundChange struct {
	Order    *models.Order
	RefundId string
}

func NewOrderRefundChangeEvent(order *models.Order, refundId string) *OrderRefundChange {
	return &OrderRefundChange{Order: order, RefundId: refundId}
}

func (o OrderRefundChange) Delay() time.Duration {
	return 0
}

func (o OrderRefundChange) Body() []byte {
	marshal, err := json.Marshal(o)
	if err != nil {
		log.Errorf("OrderRefundChange json marshal error:%s", err)
	}
	return marshal
}
