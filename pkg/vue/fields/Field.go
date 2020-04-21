package fields

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/element"
	"reflect"
	"strings"
)

type FieldRule struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func (f *FieldRule) GetRule() string {
	return f.Rule
}

func (f *FieldRule) GetMessage() string {
	return f.Message
}

type Field struct {
	*element.Element
	Panel             string                `json:"panel"`
	Readonly          bool                  `json:"readonly"`
	Expand            bool                  `json:"expand"`
	showOnIndex       bool                  `json:"-"`
	showOnDetail      bool                  `json:"-"`
	showOnCreation    bool                  `json:"-"`
	showOnUpdate      bool                  `json:"-"`
	fieldName         string                `json:"-"`
	Name              string                `json:"name"`
	Attribute         string                `json:"attribute"`
	Value             interface{}           `json:"value"`
	Sortable          bool                  `json:"sortable"`
	Nullable          bool                  `json:"nullable"`
	NullValue         interface{}           `json:"null_value"`
	TextAlign         string                `json:"text_align"`
	Stacked           bool                  `json:"stacked"`
	AsHtml            bool                  `json:"as_html"`
	Placeholder       string                `json:"placeholder"`
	Rules             []contracts.FieldRule `json:"rules"`
	resolveForDisplay func(ctx *gin.Context, model interface{}) interface{}
}

func (this Field) Call(model interface{}) {

}

func (this Field) GetRules() []contracts.FieldRule {
	return this.Rules
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

func NewField(name string, fieldName string, opts ...FieldOption) *Field {
	field := &Field{Name: name, fieldName: fieldName, Attribute: setAttribute(fieldName), Element: element.NewElement()}
	for _, opt := range opts {
		opt(field)
	}
	return field
}

// 列表页是否可见
func (this Field) ShowOnIndex() bool {
	return this.showOnIndex
}

// 详情页是否可见
func (this Field) ShowOnDetail() bool {
	return this.showOnDetail
}

// 创建页是否可见
func (this Field) ShowOnCreation() bool {
	return this.showOnCreation
}

// 更新页是否可见
func (this Field) ShowOnUpdate() bool {
	return this.showOnUpdate
}

// 设置为表格扩展行
func (f *Field) SetExpand(ok bool) {
	f.Expand = ok
	if f.Expand {
		f.showOnIndex = ok
	}
}

// 设置首页可见
func (f *Field) SetShowOnIndex(show bool) {
	f.showOnIndex = show
}

// 设置详情可见
func (f *Field) SetShowOnDetail(show bool) {
	f.showOnDetail = show
}

// 设置更新可见
func (f *Field) SetShowOnUpdate(show bool) {
	f.showOnUpdate = show
}

// 设置创建可见
func (f *Field) SetShowOnCreation(show bool) {
	f.showOnCreation = show
}

func (f *Field) HideFromIndex(cb func() bool) {
	f.showOnIndex = !cb()
}

func (f *Field) HideFromDetail(cb func() bool) {
	f.showOnDetail = !cb()
}

// 仅列表页可见
func (f *Field) OnlyOnIndex() {
	f.showOnIndex = true
	f.showOnCreation = false
	f.showOnDetail = false
	f.showOnUpdate = false
}

// 仅详情页可见
func (f *Field) OnlyOnDetail() {
	f.showOnIndex = false
	f.showOnCreation = false
	f.showOnDetail = true
	f.showOnUpdate = false
}

// 仅表单页可见
func (f *Field) OnlyOnForm() {
	f.showOnIndex = false
	f.showOnCreation = true
	f.showOnDetail = false
	f.showOnUpdate = true
}

// 除表单页外都可见
func (f *Field) ExceptOnForms() {
	f.showOnIndex = true
	f.showOnCreation = false
	f.showOnDetail = true
	f.showOnUpdate = false
}

// 计算属性，设置值
func (this *Field) ResolveForDisplay(cb func(ctx *gin.Context, model interface{}) interface{}) {
	this.resolveForDisplay = cb
}

// 赋值
func (this *Field) Resolve(ctx *gin.Context, model interface{}) {
	if this.resolveForDisplay != nil {
		this.Value = this.resolveForDisplay(ctx, model)
		return
	}
	this.Value = this.resolveAttribute(ctx, model)
}

func (this Field) resolveAttribute(ctx *gin.Context, model interface{}) interface{} {
	value := getValueByField(model, this.fieldName)
	if value == nil {
		return this.NullValue
	}
	return value

}

func (this *Field) SetPanel(name string) {
	this.Panel = name
}

func (this Field) GetPanel() string {
	return this.Panel
}

func (this Field) GetAttribute() string {
	return this.Attribute
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

	if value.IsValid() {
		f := value.FieldByName(field)

		if f.IsValid() {
			return f.Interface()
		}
	}

	return nil
}
