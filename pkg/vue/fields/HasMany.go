package fields

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/helper"
	"go-shop-v2/pkg/vue/panels"
)

// 关联集合字段
type HasMany struct {
	*Field          `inline`
	LocalKey        string             `json:"local_key"`
	resource        contracts.Resource `json:"-"`
	ViaResourceName string             `json:"via_resource_name"`
	panel           contracts.Panel    `json:"-"`
	Headings        []contracts.Field  `json:"headings"`
}

func NewHasManyField(name string, resource contracts.Resource, opts ...FieldOption) *HasMany {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(false),
		SetShowOnDetail(true),
		SetShowOnCreation(false),
		SetShowOnUpdate(false),
		WithComponent("has-many-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)

	panel := panels.NewPanel(name).SetWithoutPending(true)

	return &HasMany{
		Field:           NewField(name, "", options...),
		resource:        resource,
		ViaResourceName: helper.ResourceUriKey(resource),
		panel:           panel,
	}
}

// panel包裹
func (this *HasMany) WarpPanel() contracts.Panel {
	return this.panel
}

// 设置关联字段
func (this *HasMany) SetLocalKey(key string) *HasMany {
	this.LocalKey = key
	return this
}

// 赋值
func (this *HasMany) Resolve(ctx *gin.Context, model interface{}) {
	if this.LocalKey == "" {
		this.LocalKey = fmt.Sprintf("%s_id", utils.StrToSnake(utils.StructToName(model)))
	}
	this.Headings = helper.ResolveIndexFields(ctx, this.resource)
}
