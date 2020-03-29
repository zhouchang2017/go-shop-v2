package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

type AdminCreated struct {
	Admin *models.Admin
}

func (this AdminCreated) Body() []byte {
	bytes, _ := json.Marshal(this.Admin)
	return bytes
}

func (AdminCreated) Delay() time.Duration {
	return time.Second * 0
}
