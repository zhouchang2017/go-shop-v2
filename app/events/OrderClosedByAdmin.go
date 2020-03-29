package events

import (
	"encoding/json"
	"go-shop-v2/app/models"
	"time"
)

// 后台管理员关闭订单事件
type OrderClosedByAdmin struct {
	order *models.Order
}

func NewOrderClosedByAdminEvent(order *models.Order) *OrderClosedByAdmin {
	return &OrderClosedByAdmin{order: order}
}

func (o OrderClosedByAdmin) Body() []byte {
	bytes, _ := json.Marshal(o.order)
	return bytes
}

func (o OrderClosedByAdmin) Delay() time.Duration {
	return 0
}
