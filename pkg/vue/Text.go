package vue

type Text struct {
	*BasicField
}

func NewTextField(name string, fieldName string, opts ...FieldOption) *Text {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetComponent("text-field"),
		SetTextAlign("left"),
	}
	options = append(options,opts...)

	return &Text{BasicField: NewField(name, fieldName, options...)}
}
