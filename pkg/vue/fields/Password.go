package fields

import "github.com/gin-gonic/gin"

type Password struct {
	*Field
}

func NewPasswordField(name string, fieldName string, opts ...FieldOption) *Password {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetComponent("password-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)

	return &Password{Field: NewField(name, fieldName, options...)}
}

// Resolve the field's value.
func (this *Password) Resolve(ctx *gin.Context, model interface{}) {
	this.Value = nil
}
