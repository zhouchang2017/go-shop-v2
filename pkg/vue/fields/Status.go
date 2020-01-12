package fields

type StatusField struct {
	*Field
	Options []*StatusOption `json:"options"`
}

func NewStatusField(name string, fieldName string, opts ...FieldOption) *StatusField {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		WithComponent("status-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)

	return &StatusField{Field: NewField(name, fieldName, options...)}
}

func (this *StatusField) WithOptions(opts []*StatusOption) *StatusField {
	this.Options = opts
	return this
}

type StatusOption struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

func NewStatusOption(label string, value interface{}) *StatusOption {
	return &StatusOption{Label: label, Value: value}
}

func (s *StatusOption) SetType(statusType string) *StatusOption {
	s.Type = statusType
	return s
}

func (s *StatusOption) Success() *StatusOption {
	s.Type = "bg-green-400"
	return s
}

func (s *StatusOption) Error() *StatusOption {
	s.Type = "bg-red-400"
	return s
}

func (s *StatusOption) Info() *StatusOption {
	s.Type = "bg-blue-400"
	return s
}

func (s *StatusOption) Warning() *StatusOption {
	s.Type = "bg-yellow-400"
	return s
}

func (s *StatusOption) Danger() *StatusOption {
	s.Type = "bg-red-400"
	return s
}

func (s *StatusOption) Cancel() *StatusOption {
	s.Type = "bg-gray-400"
	return s
}
