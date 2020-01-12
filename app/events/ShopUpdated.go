package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

type ShopUpdated struct {
	Shop *models.Shop
}

func (ShopUpdated) ExchangeName() string {
	return "shop.updated"
}

func (ShopUpdated) ExchangeType() string {
	return "fanout"
}

func (ShopUpdated) RoutingKey() string {
	return "shop.updated"
}

func (this ShopUpdated) Body() []byte {
	bytes, _ := json.Marshal(this.Shop)
	return bytes
}

func (ShopUpdated) Delay() time.Duration {
	return time.Second * 0
}
