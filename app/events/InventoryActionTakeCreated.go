package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

// 出货单创建事件
type InventoryActionTakeCreated struct {
	Action *models.ManualInventoryAction
}

func (InventoryActionTakeCreated) ExchangeName() string {
	return "inventory.action.take.created"
}

func (InventoryActionTakeCreated) ExchangeType() string {
	return "fanout"
}

func (InventoryActionTakeCreated) RoutingKey() string {
	return "inventory.action.created"
}

func (this InventoryActionTakeCreated) Body() []byte {
	bytes, _ := json.Marshal(this.Action)
	return bytes
}

// 延迟触发
func (InventoryActionTakeCreated) Delay() time.Duration {
	return time.Minute * 5
}
