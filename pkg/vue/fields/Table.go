package fields

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/vue/contracts"
	"reflect"
)

type Table struct {
	*Field        `inline`
	Headings      []contracts.Field `json:"headings"`
	fieldsFactory func() []contracts.Field
}

func NewTable(name string, fieldName string, fields func() []contracts.Field, opts ...FieldOption) *Table {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(false),
		SetShowOnDetail(true),
		WithComponent("table-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)

	table := &Table{Field: NewField(name, fieldName, options...), Headings: fields()}
	table.fieldsFactory = fields
	//table.WithMeta("headings", table.Headings)
	return table
}

//
func (this *Table) Fields(fields []contracts.Field) *Table {
	this.Headings = fields
	this.fieldsFactory = func() []contracts.Field {
		return fields
	}
	return this
}

func (this *Table) resolveField(ctx *gin.Context, value interface{}) []contracts.Field {
	row := this.fieldsFactory()
	for _, field := range row {
		field.Resolve(ctx, value)
	}
	return row
}

func (this *Table) makeFields(ctx *gin.Context, value interface{}) interface{} {
	// 验证value是否是数组或者结构体
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		values := [][]contracts.Field{}
		valueOf := reflect.ValueOf(value)
		len := valueOf.Len()
		for i := 0; i < len; i++ {
			values = append(values, this.resolveField(ctx, valueOf.Index(i).Interface()))
		}

		return values
	}
	return this.resolveField(ctx, value)
}

// Resolve the field's value.
func (this *Table) Resolve(ctx *gin.Context, model interface{}) {
	if this.resolveForDisplay != nil {
		this.Value = this.makeFields(ctx, this.resolveForDisplay(ctx, model))
		return
	}
	this.Value = this.makeFields(ctx, this.resolveAttribute(ctx, model))
}
