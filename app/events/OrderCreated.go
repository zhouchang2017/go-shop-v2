package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

// 新订单事件
type OrderCreated struct {
	Order *models.Order
}

func NewOrderCreatedEvent(order *models.Order) *OrderCreated {
	return &OrderCreated{Order: order}
}

func (o OrderCreated) Body() []byte {
	bytes, _ := json.Marshal(o.Order)
	return bytes
}

func (o OrderCreated) Delay() time.Duration {
	return time.Second * 0
}
