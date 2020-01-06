package fields

type Text struct {
	*Field
}

func NewTextField(name string, fieldName string, opts ...FieldOption) *Text {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("text-field"),
		SetTextAlign("left"),
	}
	options = append(options,opts...)

	return &Text{Field: NewField(name, fieldName, options...)}
}
