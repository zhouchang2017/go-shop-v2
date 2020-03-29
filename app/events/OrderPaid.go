package events

import (
	"time"
)

// 订单已付款事件
type OrderPaid struct {
	orderOn string // 订单号
}

func NewOrderPaidEvent(orderOn string) *OrderPaid {
	return &OrderPaid{orderOn: orderOn}
}

func (o OrderPaid) Body() []byte {
	return []byte(o.orderOn)
}

func (o OrderPaid) Delay() time.Duration {
	return time.Second * 0
}
