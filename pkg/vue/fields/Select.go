package fields

type Select struct {
	*Field      `inline`
	Options     []*SelectOption `json:"options"`
}

func NewSelect(name string, fieldName string, opts ...FieldOption) *Select {
	var fieldOptions = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("select-field"),
		SetTextAlign("left"),
	}
	fieldOptions = append(fieldOptions, opts...)
	return &Select{
		Field: NewField(name, fieldName, fieldOptions...),
	}
}

func (this *Select) WithOptions(opts []*SelectOption) *Select {
	this.Options = opts
	return this
}

// 列表页调用vue组件名称
func (this *Select) IndexComponent() {
	this.ElementComponent = "text-field"
}

// 详情页调用vue组件名称
func (this *Select) DetailComponent() {
	this.ElementComponent = "text-field"
}

type SelectOption struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
	Group string      `json:"group"`
}
