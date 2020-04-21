package events

import (
	"time"
)

type OrderTimeOut struct {
	orderOn string // 订单号
}

func NewOrderTimeOutEvent(orderOn string) *OrderTimeOut {
	return &OrderTimeOut{orderOn: orderOn}
}

func (o OrderTimeOut) Body() []byte {
	return []byte(o.orderOn)
}

func (o OrderTimeOut) Delay() time.Duration {
	return time.Minute * 30
}
