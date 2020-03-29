package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

type ShopUpdated struct {
	Shop *models.Shop
}


func (this ShopUpdated) Body() []byte {
	bytes, _ := json.Marshal(this.Shop)
	return bytes
}

func (ShopUpdated) Delay() time.Duration {
	return time.Second * 0
}
