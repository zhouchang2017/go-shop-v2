package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/vue/contracts"
)

// 资源URI KEY
func ResourceUriKey(resource contracts.Resource) string {
	if customUri, ok := resource.(contracts.CustomUri); ok {
		return customUri.UriKey()
	} else {
		return utils.StructNameToSnakeAndPlural(resource)
	}
}

// 资源URI PARAM
func ResourceIdParam(resource contracts.Resource) string {
	return utils.StrToSingular(utils.StructToName(resource))
}

// 列表页字段
func ResolveIndexFields(ctx *gin.Context, resource contracts.Resource) []contracts.Field {
	fields := []contracts.Field{}
	for _, field := range resource.Fields(ctx, nil)() {
		if isField, ok := field.(contracts.Field); ok {

			if isField.ShowOnIndex() && isField.AuthorizedTo(ctx,ctx2.GetUser(ctx).(auth.Authenticatable)) {
				// 自定义列表页组件
				if hasIndexComponent, ok := isField.(contracts.CustomIndexFieldComponent); ok {
					hasIndexComponent.IndexComponent()
				}
				fields = append(fields, isField)
				continue
			}
		}

		if isPanel, ok := field.(contracts.Panel); ok {
			for _, panelField := range isPanel.GetFields() {
				if panelField.ShowOnIndex() && panelField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
					// 自定义列表页组件
					if hasIndexComponent, ok := panelField.(contracts.CustomIndexFieldComponent); ok {
						hasIndexComponent.IndexComponent()
					}
					fields = append(fields, panelField)
				}
			}
		}
	}
	return fields
}

// 聚合URI KEY
func LensUriKey(lens contracts.Lens) string {
	if customUri, ok := lens.(contracts.CustomUri); ok {
		return customUri.UriKey()
	} else {
		return utils.StructNameToSnakeAndPlural(lens)
	}
}

// vue router
// 列表页路由名称
func IndexRouteName(resource contracts.Resource) string {
	if implement, ok := resource.(contracts.ResourceCustomIndexComponent); ok {
		return implement.IndexComponent().VueRouter().RouterName()
	}
	return fmt.Sprintf("%s.index", ResourceUriKey(resource))
}

// 详情页路由名称
func DetailRouteName(resource contracts.Resource) string {
	if implement, ok := resource.(contracts.ResourceCustomDetailComponent); ok {
		return implement.DetailComponent().VueRouter().RouterName()
	}
	return fmt.Sprintf("%s.detail", ResourceUriKey(resource))
}

// 更新页路由名称
func UpdateRouteName(resource contracts.Resource) string {
	if implement, ok := resource.(contracts.ResourceCustomUpdateComponent); ok {
		return implement.UpdateComponent().VueRouter().RouterName()
	}
	return fmt.Sprintf("%s.edit", ResourceUriKey(resource))
}

// 创建页路由名称
func CreationRouteName(resource contracts.Resource) string {
	if implement, ok := resource.(contracts.ResourceCustomCreationComponent); ok {
		return implement.CreationComponent().VueRouter().RouterName()
	}
	return fmt.Sprintf("%s.create", ResourceUriKey(resource))
}

// 聚合页路由名称Lens
func LensRouteName(resource contracts.Resource, lens contracts.Lens) string {
	return fmt.Sprintf("%s.lenses.%s", ResourceUriKey(resource), LensUriKey(lens))
}