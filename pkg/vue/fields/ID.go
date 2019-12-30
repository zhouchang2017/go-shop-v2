package fields

// id字段
type ID struct {
	*Field `inline`
}

func NewIDField(opts ...FieldOption) *ID {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetComponent("text-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)

	return &ID{Field: NewField("ID", "ID", options...)}
}
