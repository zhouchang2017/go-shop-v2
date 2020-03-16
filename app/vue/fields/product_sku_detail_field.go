package fields

import (
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

// 自定义 产品详情页sku列表字段
type ProductSkuDetailField struct {
	*fields.Field `inline`
	panel         *panels.Panel
}

func (this ProductSkuDetailField) WarpPanel() contracts.Panel {
	return this.panel
}

func NewProductSkuDetailField(name string, fieldName string, opts ...fields.FieldOption) *ProductSkuDetailField {
	var fieldOptions = []fields.FieldOption{
		fields.SetPrefixComponent(true),
		fields.SetShowOnIndex(false),
		fields.SetShowOnDetail(true),
		fields.SetShowOnCreation(false),
		fields.SetShowOnUpdate(false),
		fields.WithComponent("product-sku-field"),
		fields.SetTextAlign("left"),
	}

	panel := panels.NewPanel(name).SetWithoutPending(true)

	return &ProductSkuDetailField{
		Field: fields.NewField(name, fieldName, fieldOptions...),
		panel: panel,
	}
}
