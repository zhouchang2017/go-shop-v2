package vue

type DateTime struct {
	*BasicField
}

func NewDateTime(name string, fieldName string,opts ...FieldOption) *DateTime {
	opts = append(opts,
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetComponent("date-time-field"),
		SetTextAlign("left"),
	)

	return &DateTime{BasicField:NewField(name,fieldName,opts...)}
}
