package fields

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

type OrderItemsField struct {
	*fields.Field `inline`
	panel         *panels.Panel
}

func (this OrderItemsField) WarpPanel() contracts.Panel {
	return this.panel
}

func NewOrderItemsField(opts ...fields.FieldOption) *OrderItemsField {
	var fieldOptions = []fields.FieldOption{
		fields.SetPrefixComponent(true),
		fields.SetShowOnIndex(false),
		fields.SetShowOnDetail(true),
		fields.SetShowOnCreation(false),
		fields.SetShowOnUpdate(false),
		fields.WithComponent("order-items-field"),
		fields.SetTextAlign("left"),
	}

	panel := panels.NewPanel("订单信息").SetWithoutPending(true)

	return &OrderItemsField{
		Field: fields.NewField("订单状态", "status", fieldOptions...),
		panel: panel,
	}
}

// 赋值
func (this *OrderItemsField) Resolve(ctx *gin.Context, model interface{}) {
	this.Value = model
}
