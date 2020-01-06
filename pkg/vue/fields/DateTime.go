package fields

type DateTime struct {
	*Field
}

func NewDateTime(name string, fieldName string, opts ...FieldOption) *DateTime {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		WithComponent("date-time-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)

	return &DateTime{Field: NewField(name, fieldName, options...)}
}
