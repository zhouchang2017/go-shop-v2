package vue

import (
	"fmt"
	"go-shop-v2/pkg/utils"
)

type ResourceHelper struct {
	resource Resource
}

func NewResourceHelper(resource Resource) *ResourceHelper {
	return &ResourceHelper{resource: resource}
}

// Get the displayable label of the resource.
func (this *ResourceHelper) Label() string {
	return utils.StrToPlural(utils.StructToName(this.resource))
}

// Get the displayable singular label of the resource.
func (this *ResourceHelper) SingularLabel() string {
	return utils.StrToSingular(utils.StructToName(this.resource))
}


func (this *ResourceHelper) ResourceID() string {
	return utils.StrToSingular(utils.StructToName(this.resource))
}

func (this *ResourceHelper) VueUriKey() string {
	if custom, ok := this.resource.(CustomVueUriKey); ok {
		return custom.CustomVueUriKey()
	}
	return this.UriKey()
}

// Get the URI key for the resource.
func (this *ResourceHelper) UriKey() string {
	return utils.StructNameToSnakeAndPlural(this.resource)
}

// vue列表路由
func (this *ResourceHelper) IndexRouterName() string {
	return fmt.Sprintf("%s.index", this.UriKey())
}

// vue详情路由
func (this *ResourceHelper) DetailRouterName() string {
	return fmt.Sprintf("%s.detail", this.UriKey())
}

// vue编辑路由
func (this *ResourceHelper) EditRouterName() string {
	return fmt.Sprintf("%s.edit", this.UriKey())
}

// vue创建路由
func (this *ResourceHelper) CreateRouterName() string {
	return fmt.Sprintf("%s.create", this.UriKey())
}
