package fields

import "go-shop-v2/pkg/vue/contracts"

type Select struct {
	*Field      `inline`
	optionField contracts.Field
	Options     []contracts.Field `json:"options"`
}

func NewSelect(name string, fieldName string, opts ...FieldOption) *Select {
	var fieldOptions = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		SetComponent("select-field"),
		SetTextAlign("left"),
	}
	fieldOptions = append(fieldOptions, opts...)
	return &Select{
		Field:       NewField(name, fieldName, fieldOptions...),
	}
}

func (this *Select) SetOptions(opts []contracts.Field) *Select {
	//this.optionField
	this.Options = opts
	return this
}

// 列表页调用vue组件名称
func (this *Select) IndexComponent() {
	this.Component = "text-field"
}

// 详情页调用vue组件名称
func (this *Select) DetailComponent() {
	this.Component = "text-field"
}

type SelectOption struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
	Group string      `json:"group"`
	Field contracts.Field
}
