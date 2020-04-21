package core

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/panels"
)

// 列表页字段
func resolveLensIndexFields(ctx *gin.Context, lens contracts.Lens) []contracts.Field {
	fields := []contracts.Field{}
	for _, field := range lens.Fields(ctx, nil)() {
		if isField, ok := field.(contracts.Field); ok {

			if isField.ShowOnIndex() && isField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				// 自定义列表页组件
				if hasIndexComponent, ok := isField.(contracts.CustomIndexFieldComponent); ok {
					hasIndexComponent.IndexComponent()
				}
				fields = append(fields, isField)
				continue
			}
		}

		if isPanel, ok := field.(*panels.Panel); ok {
			for _, panelField := range isPanel.Fields {
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

// 聚合列表页字段格式化输出
func SerializeLens(ctx *gin.Context, lens contracts.Lens, model interface{}) map[string]interface{} {
	var maps = map[string]interface{}{}
	fields := []contracts.Field{}
	for _, field := range resolveLensIndexFields(ctx, lens) {
		field.Resolve(ctx, model)
		fields = append(fields, field)

		field.Call(model)
	}
	maps["fields"] = fields
	return maps
}
