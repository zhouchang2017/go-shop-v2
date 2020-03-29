package events

import "time"

// 买家取消订单事件
type OrderClosedByUser struct {
	orderOn string // 订单号
}

func NewOrderClosedByUserEvent(orderOn string) *OrderClosedByUser {
	return &OrderClosedByUser{orderOn: orderOn}
}

func (o OrderClosedByUser) Body() []byte {
	return []byte(o.orderOn)
}

func (o OrderClosedByUser) Delay() time.Duration {
	return time.Second * 0
}
 
