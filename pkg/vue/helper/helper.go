package helper

import (
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
