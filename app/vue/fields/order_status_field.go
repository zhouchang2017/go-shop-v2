package fields

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

type OrderStatusField struct {
	*fields.Field `inline`
	panel         *panels.Panel
}

func (this OrderStatusField) WarpPanel() contracts.Panel {
	return this.panel
}

func NewOrderStatusField(opts ...fields.FieldOption) *OrderStatusField {
	var fieldOptions = []fields.FieldOption{
		fields.SetPrefixComponent(true),
		fields.SetShowOnIndex(true),
		fields.SetShowOnDetail(true),
		fields.SetShowOnCreation(false),
		fields.SetShowOnUpdate(false),
		fields.WithComponent("order-status-field"),
		fields.SetTextAlign("left"),
	}

	panel := panels.NewPanel("订单状态").SetWithoutPending(true).SetSort(100)

	return &OrderStatusField{
		Field: fields.NewField("订单状态", "status", fieldOptions...),
		panel: panel,
	}
}

// 赋值
func (this *OrderStatusField) Resolve(ctx *gin.Context, model interface{}) {
	this.Value = model
}
