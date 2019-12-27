package vue

type Panel struct {
	BasicElement
	Name   string  `json:"name"`
	Limit  *int64  `json:"limit"`
	Fields []Field `json:"-"`
}

func NewPanel(name string, fields ...Field) *Panel {
	p := &Panel{Name: name}
	p.PrepareFields(fields...)
	return p
}

func (p *Panel) PrepareFields(fields ...Field) {
	for _, field := range fields {
		field.SetPanel(p.Name)
		p.Fields = append(p.Fields, field)
	}
}

// Set the number of initially visible fields.
func (p *Panel) SetLimit(num int64) {
	p.Limit = &num
}
