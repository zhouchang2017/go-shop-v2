package events

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"time"
)

type OrderRefundCancelByUser struct {
	refund *models.Refund
}

func NewOrderRefundCancelByUserEvent(refund *models.Refund) *OrderRefundCancelByUser {
	return &OrderRefundCancelByUser{refund}
}

func (o OrderRefundCancelByUser) Delay() time.Duration {
	return 0
}

func (o OrderRefundCancelByUser) Body() []byte {
	marshal, err := json.Marshal(o.refund)
	if err != nil {
		log.Errorf("OrderRefundChange json marshal error:%s", err)
	}
	return marshal
}
