package fields

type CheckboxGroup struct {
	*Field
	CheckboxGroupOptions []*CheckboxGroupOption `json:"options"`
	cbOption             func() []*CheckboxGroupOption
	KeyField             string `json:"key_field"`
}

type CheckboxGroupOption struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

func NewCheckboxGroup(name string, fieldName string, opts ...FieldOption) *CheckboxGroup {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(false),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("checkbox-group-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)

	return &CheckboxGroup{Field: NewField(name, fieldName, options...)}
}

func (this *CheckboxGroup) AsyncLoadOptions() {

}

func (this *CheckboxGroup) Call() {
	if this.cbOption != nil {
		this.CheckboxGroupOptions = this.cbOption()
	}
}

func (this *CheckboxGroup) CallbackOptions(cb func() []*CheckboxGroupOption) *CheckboxGroup {
	this.cbOption = cb
	return this
}

func (this *CheckboxGroup) Key(key string) *CheckboxGroup {
	this.KeyField = key
	return this
}

func (this *CheckboxGroup) Options(opts ...*CheckboxGroupOption) *CheckboxGroup {
	this.CheckboxGroupOptions = opts
	return this
}
