package events

import (
	"time"
)

type OrderRefundDone struct {
	orderId string
}

func NewOrderRefundDoneEvent(orderId string) *OrderRefundDone {
	return &OrderRefundDone{orderId: orderId}
}

func (o OrderRefundDone) Delay() time.Duration {
	return 0
}

func (o OrderRefundDone) Body() []byte {
	return []byte(o.orderId)
}
