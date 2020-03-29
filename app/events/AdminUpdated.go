package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

type AdminUpdated struct {
	Admin *models.Admin
}


func (this AdminUpdated) Body() []byte {
	bytes, _ := json.Marshal(this.Admin)
	return bytes
}

func (AdminUpdated) Delay() time.Duration {
	return time.Second * 0
}
