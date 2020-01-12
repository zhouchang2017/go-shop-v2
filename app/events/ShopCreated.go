package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

type ShopCreated struct {
	Shop *models.Shop
}

func (ShopCreated) ExchangeName() string {
	return "shop.created"
}

func (ShopCreated) ExchangeType() string {
	return "fanout"
}

func (ShopCreated) RoutingKey() string {
	return "shop.created"
}

func (this ShopCreated) Body() []byte {
	bytes, _ := json.Marshal(this.Shop)
	return bytes
}

func (ShopCreated) Delay() time.Duration {
	return time.Second * 0
}
