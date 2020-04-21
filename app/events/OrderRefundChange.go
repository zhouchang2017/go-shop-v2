package events

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"time"
)

type OrderRefundChange struct {
	Refund *models.Refund
}

func NewOrderRefundChangeEvent(refund *models.Refund) *OrderRefundChange {
	return &OrderRefundChange{refund}
}

func (o OrderRefundChange) Delay() time.Duration {
	return 0
}

func (o OrderRefundChange) Body() []byte {
	marshal, err := json.Marshal(o.Refund)
	if err != nil {
		log.Errorf("OrderRefundChange json marshal error:%s", err)
	}
	return marshal
}
