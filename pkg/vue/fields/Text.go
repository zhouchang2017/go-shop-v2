package fields

import (
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/helper"
)

type Text struct {
	*Field
	LinkOption   map[string]interface{} `json:"link_option"`
	LinkResource contracts.Resource
	Pk           string
}

// https://element.eleme.cn/#/zh-CN/component/input
func NewTextField(name string, fieldName string, opts ...FieldOption) *Text {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("text-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)

	return &Text{Field: NewField(name, fieldName, options...)}
}

func (this *Text) Text() *Text {
	this.WithMeta("type", "text")
	return this
}

func (this *Text) Textarea() *Text {
	this.SetShowOnIndex(false)
	this.WithMeta("type", "textarea")
	return this
}

func (this *Text) Email() *Text {
	this.WithMeta("type", "email")
	return this
}

func (this *Text) ShowWordLimit() *Text {
	this.WithMeta("showWordLimit", true)
	return this
}

func (this *Text) Clearable() *Text {
	this.WithMeta("clearable", true)
	return this
}

func (this *Text) Rows(num int64) *Text {
	this.WithMeta("rows", num)
	return this
}

func (this *Text) Min(num int64) *Text {
	this.WithMeta("min", num)
	return this
}

func (this *Text) Max(num int64) *Text {
	this.WithMeta("max", num)
	return this
}

func (this *Text) Step(num int64) *Text {
	this.WithMeta("step", num)
	return this
}

func (this *Text) Pattern(reg string) *Text {
	this.WithMeta("pattern", reg)
	return this
}

func (this *Text) InputNumber() *Text {
	this.WithMeta("input_number", true)
	this.Value = 0
	return this
}

type AutosizeOpt func(text *Text)

func MinRows(num int64) AutosizeOpt {
	return func(text *Text) {
		text.WithMeta("minRows", num)
	}
}

func MaxRows(num int64) AutosizeOpt {
	return func(text *Text) {
		text.WithMeta("maxRows", num)
	}
}

func (this *Text) Autosize(opts ...AutosizeOpt) *Text {
	for _, opt := range opts {
		opt(this)
	}
	this.WithMeta("autosize", true)
	return this
}

func (this *Text) Link(resource contracts.Resource, pk string) *Text {

	if pk == "" {
		pk = this.Attribute
	}
	this.Pk = pk
	this.LinkResource = resource

	return this
}

func (this *Text) Call(model interface{}) {
	if this.LinkResource != nil {
		name := helper.DetailRouteName(this.LinkResource)
		if model != nil {
			id := getValueByField(model, this.Pk)
			this.LinkOption = map[string]interface{}{
				"routeName": name,
				"id":        id,
			}
		}

	}
}
