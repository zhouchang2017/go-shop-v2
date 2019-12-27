package vue

// id字段
type ID struct {
	*BasicField `inline`
}

func NewIDField(opts ...FieldOption) *ID {
	opts = append(opts,
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetComponent("text-field"),
		SetTextAlign("left"),
	)
	return &ID{BasicField: NewField("ID", "ID", opts...)}
}
