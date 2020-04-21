package events

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"time"
)

// 订单申请退款
type OrderApplyRefund struct {
	Refund *models.Refund `json:"refund"`
}

func (o OrderApplyRefund) Delay() time.Duration {
	return 0
}

func (o OrderApplyRefund) Body() []byte {
	marshal, err := json.Marshal(o.Refund)
	if err != nil {
		log.Errorf("OrderApplyRefund json marshal error:%s", err)
	}
	return marshal
}

func NewOrderApplyRefundEvent(refund *models.Refund) *OrderApplyRefund {
	return &OrderApplyRefund{refund}
}
