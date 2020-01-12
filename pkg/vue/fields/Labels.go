package fields

// 标签字段，值为一个数组
type LabelsFields struct {
	*Field `inline`
}

func NewLabelsFields(name string, fieldName string, opts ...FieldOption) *LabelsFields {
	var fieldOptions = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(false),
		SetShowOnUpdate(false),
		WithComponent("labels-field"),
		SetTextAlign("left"),
	}
	fieldOptions = append(fieldOptions, opts...)
	return &LabelsFields{
		Field: NewField(name, fieldName, fieldOptions...),
	}
}

// 显示文字的字段
func (this *LabelsFields) Label(label string) *LabelsFields {
	this.WithMeta("label", label)
	return this
}

// 显示鼠标经过 tooltip 字段，如果没有就不显示
func (this *LabelsFields) Tooltip(field string) *LabelsFields {
	this.WithMeta("tooltip", field)
	return this
}
