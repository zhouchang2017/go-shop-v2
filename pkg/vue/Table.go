package vue

type Table struct {
	*BasicField `inline`
	Headings    map[string]string `json:"-"`
}

func NewTable(name string, fieldName string, heading map[string]string, opts ...FieldOption) *Table {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(false),
		SetShowOnDetail(true),
		SetComponent("table-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)

	table := &Table{BasicField: NewField(name, fieldName, options...), Headings: heading}
	table.WithMeta("headings", table.Headings)
	return table
}
