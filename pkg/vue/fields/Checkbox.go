package fields

type Checkbox struct {
	*Field
}

func NewCheckbox(name string, fieldName string, opts ...FieldOption) *Checkbox {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetComponent("checkbox-field"),
		SetTextAlign("left"),
	}
	options = append(options,opts...)

	return &Checkbox{Field: NewField(name, fieldName, options...)}
} 
