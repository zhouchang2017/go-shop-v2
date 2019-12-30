package vue

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/utils"
	"reflect"
	"strings"
)

type HasFields interface {
	Fields(ctx *gin.Context, model interface{}) func() []interface{}
}

type Field interface {
	Element
	ShowOnIndex() bool
	ShowOnDetail() bool
	ShowOnCreation() bool
	ShowOnUpdate() bool
	Resolve(ctx *gin.Context, model interface{})
	SetPanel(name string)
	GetPanel() string
	GetRules() []*FieldRule
	GetAttribute() string
	Fill(ctx *gin.Context, data map[string]interface{}, model interface{})
}

type FieldOption func(field interface{})

func SetSortable(sort bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.SetSortable(sort)
		}
	}
}

func SetRules(rules []*FieldRule) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.Rules = rules
		}
	}
}

func SetAttribute(attr string) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.SetAttribute(attr)
		}
	}
}

func SetNullable(nullable bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.SetNullable(nullable)
		}
	}
}

func SetNullValue(value interface{}) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.NullValue = value
		}
	}
}

func SetTextAlign(align string) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.TextAlign = align
		}
	}
}

func SetAsHtml(flag bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.AsHtml = flag
		}
	}
}

func resolveBasicField(field interface{}) (*BasicField, error) {
	if basicField, ok := field.(*BasicField); ok {
		return basicField, nil
	}
	if reflect.ValueOf(field).Kind() == reflect.Ptr {
		elem := reflect.ValueOf(field).Elem()
		for i := 0; i < elem.NumField(); i++ {
			value := elem.Field(i)
			if value.IsValid() && value.Type() == reflect.ValueOf(&BasicField{}).Type() {
				return value.Interface().(*BasicField), nil
			}
		}
	}
	return nil, fmt.Errorf("basic field not found in %+v\n", field)
}

type FieldRule struct {
	Rule    string
	Message string
}

type BasicField struct {
	FieldElement
	fieldName         string      `json:"-"`
	Name              string      `json:"name"`
	Attribute         string      `json:"attribute"`
	Value             interface{} `json:"value"`
	Sortable          bool        `json:"sortable"`
	Nullable          bool        `json:"nullable"`
	NullValue         interface{} `json:"null_value"`
	TextAlign         string      `json:"text_align"`
	Stacked           bool        `json:"stacked"`
	AsHtml            bool        `json:"as_html"`
	Rules             []*FieldRule
	resolveForDisplay func(ctx *gin.Context, model interface{}) interface{}
}

func setAttribute(fieldName string) (name string) {

	split := strings.Split(fieldName, ".")

	if len(split) > 0 {
		var names []string
		for _, s := range split {
			names = append(names, utils.StrToSnake(s))
		}
		return strings.Join(names, ".")
	}
	return utils.StrToSnake(fieldName)
}

func NewField(name string, fieldName string, opts ...FieldOption) *BasicField {
	field := &BasicField{Name: name, fieldName: fieldName, Attribute: setAttribute(fieldName)}
	for _, opt := range opts {
		opt(field)
	}
	return field
}

func (this *BasicField) SetAsHtml(flag bool) {
	this.AsHtml = flag
}

func (this *BasicField) SetSortable(sort bool) {
	this.Sortable = sort
}

func (this *BasicField) SetNullable(nullable bool) {
	this.Nullable = nullable
}

func (this *BasicField) SetTextAlign(align string) {
	this.TextAlign = align
}

func (this *BasicField) SetFieldName(field string) {
	this.fieldName = field
}

func (this *BasicField) SetAttribute(attr string) {
	this.Attribute = attr
}

// Set the help text for the field.
func (this *BasicField) Help(helpText string) {
	this.WithMeta("helpText", helpText)
}

// Resolve the field's value for display.
func (this *BasicField) ResolveForDisplay(cb func(ctx *gin.Context, model interface{}) interface{}) {
	this.resolveForDisplay = cb
}

// Resolve the field's value.
func (this *BasicField) Resolve(ctx *gin.Context, model interface{}) {
	if this.resolveForDisplay != nil {
		this.Value = this.resolveForDisplay(ctx, model)
		return
	}
	this.Value = this.ResolveAttribute(ctx, model)
}

func getValueByField(model interface{}, field string) interface{} {

	var value reflect.Value
	switch reflect.ValueOf(model).Kind() {
	case reflect.Ptr:
		value = reflect.ValueOf(model).Elem()
	case reflect.Struct:
		value = reflect.ValueOf(model)
	}
	attrs := strings.Split(field, ".")
	if len(attrs) > 1 {
		head := attrs[0]
		others := strings.Join(attrs[1:], ".")
		target := value.FieldByName(head)
		if target.IsValid() {
			return getValueByField(target.Interface(), others)
		}
	}

	f := value.FieldByName(field)
	if f.IsValid() {
		return f.Interface()
	}
	return nil
}

func (this *BasicField) GetAttribute() string {
	return this.Attribute
}

func (this *BasicField) GetRules() []*FieldRule {
	return this.Rules
}

// Resolve the given attribute from the given resource.
func (this *BasicField) ResolveAttribute(ctx *gin.Context, model interface{}) interface{} {
	return getValueByField(model, this.fieldName)
}

func (this *BasicField) Fill(ctx *gin.Context, data map[string]interface{}, model interface{}) {

	if value, ok := data[this.Attribute]; ok {
		spew.Dump(value)
	}

}
