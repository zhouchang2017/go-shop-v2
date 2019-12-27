package vue

type DateTime struct {
	*BasicField
}

func NewDateTime(name string, fieldName string,opts ...FieldOption) *DateTime {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetComponent("date-time-field"),
		SetTextAlign("left"),
	}
	options = append(options,opts...)

	return &DateTime{BasicField:NewField(name,fieldName,options...)}
}
