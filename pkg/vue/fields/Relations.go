package fields

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/helper"
	"go-shop-v2/pkg/vue/panels"
)

type Relations struct {
	*Field       `inline`
	ResourceName string `json:"resource_name"`
	panel        contracts.Panel
	Headings     []contracts.Field `json:"headings"`
	resource     contracts.ResourceRelations
	DetailRouteName string `json:"detail_route_name"`
}

func NewRelationsField(resource contracts.ResourceRelations, fieldName string, opts ...FieldOption) *Relations {
	var fieldOptions = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(false),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("relations-field"),
		SetTextAlign("left"),
	}
	fieldOptions = append(fieldOptions, opts...)
	panel := panels.NewPanel(resource.Title()).SetWithoutPending(true)
	return &Relations{
		Field:        NewField(resource.Title(), fieldName, fieldOptions...),
		ResourceName: helper.ResourceUriKey(resource),
		panel:        panel,
		resource:     resource,
		DetailRouteName:helper.DetailRouteName(resource),
	}
}

func (this *Relations) WarpPanel() contracts.Panel {
	return this.panel
}

func (this *Relations) Searchable() *Relations {
	this.WithMeta("searchable", true)
	return this
}

func (this *Relations) Multiple() *Relations {
	this.WithMeta("multiple", true)
	return this
}

func (this *Relations) MultipleLimit(num int64) *Relations {
	this.WithMeta("multiple_limit", num)
	return this
}

func (this *Relations) WithName(name string) *Relations {
	this.Name = name
	this.panel.SetName(name)
	return this
}

// 赋值
func (this *Relations) Resolve(ctx *gin.Context, model interface{}) {
	this.Field.Resolve(ctx, model)
	this.Headings = helper.ResolveIndexFields(ctx, this.resource)
}
