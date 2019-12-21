package vue

import "github.com/gin-gonic/gin"

type Field struct {
	FieldElement
	Name      string      `json:"name"`
	Attribute string      `json:"attribute"`
	Value     interface{} `json:"value"`
	Sortable  bool        `json:"sortable"`
	Nullable  bool        `json:"nullable"`
	NullValue interface{} `json:"null_value"`
	TextAlign string      `json:"text_align"`
	Stacked   bool        `json:"stacked"`
}

func NewField(name string, attribute string) *Field {
	return &Field{Name: name, Attribute: attribute}
}

// Set the help text for the field.
func (this *Field) Help(helpText string) {
	this.WithMeta("helpText", helpText)
}

// Resolve the field's value for display.
func (this *Field) ResolveForDisplay(resource interface{}) {

}

// Resolve the field's value.
func (this *Field) Resolve(resource interface{}, attr ...string) {
	if len(attr) == 0 {

	}
}

// Resolve the given attribute from the given resource.
func (this *Field) ResolveAttribute(resource interface{}, attribute string) {

}

func (this *Field) Fill(ctx *gin.Context,model interface{}) {

}
