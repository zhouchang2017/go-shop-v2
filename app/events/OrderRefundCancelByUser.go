package events

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"time"
)

type OrderRefundCancelByUser struct {
	Order    *models.Order
	RefundId string
}

func NewOrderRefundCancelByUserEvent(order *models.Order, refundId string) *OrderRefundCancelByUser {
	return &OrderRefundCancelByUser{Order: order, RefundId: refundId}
}

func (o OrderRefundCancelByUser) Delay() time.Duration {
	return 0
}

func (o OrderRefundCancelByUser) Body() []byte {
	marshal, err := json.Marshal(o)
	if err != nil {
		log.Errorf("OrderRefundChange json marshal error:%s", err)
	}
	return marshal
}
