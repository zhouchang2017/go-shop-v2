package fields

import (
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

// 自定义 产品类目销售属性 字段
type CategoryOptionField struct {
	*fields.Field `inline`
	panel         *panels.Panel
}

func (this CategoryOptionField) WarpPanel() contracts.Panel {
	return this.panel
}

func NewCategoryOptionField(name string, fieldName string, opts ...fields.FieldOption) *CategoryOptionField {
	var fieldOptions = []fields.FieldOption{
		fields.SetPrefixComponent(true),
		fields.SetShowOnIndex(false),
		fields.SetShowOnDetail(true),
		fields.SetShowOnCreation(true),
		fields.SetShowOnUpdate(true),
		fields.WithComponent("category-option-field"),
		fields.SetTextAlign("left"),
	}

	panel := panels.NewPanel(name).SetWithoutPending(true)

	return &CategoryOptionField{
		Field: fields.NewField(name, fieldName, fieldOptions...),
		panel: panel,
	}
}
