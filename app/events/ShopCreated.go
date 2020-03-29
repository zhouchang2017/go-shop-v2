package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

type ShopCreated struct {
	Shop *models.Shop
}

func (this ShopCreated) Body() []byte {
	bytes, _ := json.Marshal(this.Shop)
	return bytes
}

func (ShopCreated) Delay() time.Duration {
	return time.Second * 0
}
