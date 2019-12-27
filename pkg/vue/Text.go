package vue

type Text struct {
	*BasicField
}

func NewTextField(name string, fieldName string, opts ...FieldOption) *Text {
	opts = append(opts,
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetComponent("text-field"),
		SetTextAlign("left"),
	)

	return &Text{BasicField: NewField(name, fieldName, opts...)}
}
