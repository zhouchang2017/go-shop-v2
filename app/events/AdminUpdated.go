package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

type AdminUpdated struct {
	Admin *models.Admin
}

func (AdminUpdated) ExchangeName() string {
	return "admin.updated"
}

func (AdminUpdated) ExchangeType() string {
	return "fanout"
}

func (AdminUpdated) RoutingKey() string {
	return "admin.updated"
}

func (this AdminUpdated) Body() []byte {
	bytes, _ := json.Marshal(this.Admin)
	return bytes
}

func (AdminUpdated) Delay() time.Duration {
	return time.Second * 0
}
