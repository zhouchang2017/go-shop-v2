package events

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"time"
)

type OrderRefundDone struct {
	refund *models.Refund
}

func NewOrderRefundDoneEvent(refund *models.Refund) *OrderRefundDone {
	return &OrderRefundDone{refund}
}

func (o OrderRefundDone) Delay() time.Duration {
	return 0
}

func (o OrderRefundDone) Body() []byte {
	marshal, err := json.Marshal(o.refund)
	if err != nil {
		log.Errorf("OrderRefundDone json marshal error:%s", err)
	}
	return marshal
}
