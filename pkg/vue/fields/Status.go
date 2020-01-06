package fields

type Status struct {
	*Field
}

func NewStatusField(name string, fieldName string, opts ...FieldOption) *Status {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		WithComponent("text-field"),
		SetTextAlign("left"),
	}
	options = append(options,opts...)

	return &Status{Field: NewField(name, fieldName, options...)}
}
