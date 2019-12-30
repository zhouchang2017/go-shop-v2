package panels

import (
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/element"
)

type Panel struct {
	*element.Element
	Name           string            `json:"name"`
	Limit          *int64            `json:"limit"`
	ShowToolbar    bool              `json:"show_toolbar"`
	Fields         []contracts.Field `json:"-"`
	WithoutPending bool              `json:"without_pending"`
}

func NewPanel(name string, fields ...contracts.Field) *Panel {
	p := &Panel{Element:element.NewElement(),Name: name}
	p.SetComponent("panel")
	p.PrepareFields(fields...)
	return p
}

func (p *Panel) SetWithoutPending(flag bool) *Panel {
	p.WithoutPending = flag
	return p
}

func (p *Panel) PrepareFields(fields ...contracts.Field) {
	for _, field := range fields {
		field.SetPanel(p.Name)
		p.Fields = append(p.Fields, field)
	}
}

// Set the number of initially visible fields.
func (p *Panel) SetLimit(num int64) {
	p.Limit = &num
}
