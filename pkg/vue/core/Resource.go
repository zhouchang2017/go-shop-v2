package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

type warp struct {
	resource contracts.Resource
	httpHandler *resourceHttpHandle
	vueRouterFactory *vueRouterFactory
}

// 列表页字段
func resolveIndexFields(ctx *gin.Context, resource contracts.Resource) []contracts.Field {
	var fields []contracts.Field

	for _, field := range resource.Fields(ctx, nil)() {
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

// 资源详情页字段
func resolveDetailFields(ctx *gin.Context, resource contracts.Resource) ([]contracts.Field, []*panels.Panel) {
	var fields []contracts.Field
	var panel []*panels.Panel
	for _, field := range resource.Fields(ctx, resource.Model())() {

		if isField, ok := field.(contracts.Field); ok {
			if isField.ShowOnDetail() && isField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {

				// 自定义详情页组件
				if hasDetailComponent, ok := isField.(contracts.CustomDetailFieldComponent); ok {
					hasDetailComponent.DetailComponent()
				}

				fields = append(fields, isField)
				continue
			}
		}

		if isPanel, ok := field.(*panels.Panel); ok {
			availableFieldNum := 0
			for _, panelField := range isPanel.Fields {
				if panelField.ShowOnDetail() && panelField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
					availableFieldNum++
					// 自定义详情页组件
					if hasDetailComponent, ok := panelField.(contracts.CustomDetailFieldComponent); ok {
						hasDetailComponent.DetailComponent()
					}

					fields = append(fields, panelField)
				}
			}
			if availableFieldNum > 0 {
				panel = append(panel, isPanel)
			}

		}

	}
	return fields, panel
}

// 资源创建页字段
func resolveCreationFields(ctx *gin.Context, resource contracts.Resource) ([]contracts.Field, []*panels.Panel) {
	var fields []contracts.Field
	var panel []*panels.Panel
	defaultPanel := panels.NewPanel(fmt.Sprintf("创建%s", resource.Title()))
	panel = append(panel, defaultPanel)

	for _, field := range resource.Fields(ctx, nil)() {

		if isField, ok := field.(contracts.Field); ok {
			if isField.ShowOnCreation() && isField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {

				// 自定义创建页组件
				if hasDetailComponent, ok := isField.(contracts.CustomCreationFieldComponent); ok {
					hasDetailComponent.CreationComponent()
				}

				fields = append(fields, isField)
				continue
			}
		}

		if isPanel, ok := field.(*panels.Panel); ok {
			availableFieldNum := 0
			for _, panelField := range isPanel.Fields {
				if panelField.ShowOnCreation() && panelField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
					availableFieldNum++
					// 自定义创建页组件
					if hasDetailComponent, ok := panelField.(contracts.CustomCreationFieldComponent); ok {
						hasDetailComponent.CreationComponent()
					}

					fields = append(fields, panelField)
				}
			}
			if availableFieldNum > 0 {
				panel = append(panel, isPanel)
			}

		}

	}
	return fields, panel
}

// 资源更新页字段
func resolveUpdateFields(ctx *gin.Context, resource contracts.Resource) ([]contracts.Field, []*panels.Panel) {
	var fields []contracts.Field
	var panel []*panels.Panel
	// TODO 自定义panel title
	defaultPanel := panels.NewPanel(fmt.Sprintf("更新%s", resource.Title()))
	panel = append(panel, defaultPanel)

	for _, field := range resource.Fields(ctx, resource.Model())() {

		if isField, ok := field.(contracts.Field); ok {
			if isField.ShowOnUpdate() && isField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {

				// 自定义更新页组件
				if hasDetailComponent, ok := isField.(contracts.CustomUpdateFieldComponent); ok {
					hasDetailComponent.UpdateComponent()
				}
				fields = append(fields, isField)
				continue
			}
		}

		if isPanel, ok := field.(*panels.Panel); ok {
			availableFieldNum := 0
			for _, panelField := range isPanel.Fields {
				if panelField.ShowOnUpdate() && panelField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
					availableFieldNum++
					// 自定义更新页组件
					if hasDetailComponent, ok := panelField.(contracts.CustomUpdateFieldComponent); ok {
						hasDetailComponent.UpdateComponent()
					}
					fields = append(fields, panelField)
				}
			}
			if availableFieldNum > 0 {
				panel = append(panel, isPanel)
			}

		}

	}
	return fields, panel
}

// 列表页数据格式
func SerializeForIndex(ctx *gin.Context, resource contracts.Resource) map[string]interface{} {
	var maps = map[string]interface{}{}
	maps["AuthorizedToView"] = AuthorizedToView(ctx, resource)
	maps["AuthorizedToUpdate"] = AuthorizedToUpdate(ctx, resource)
	maps["AuthorizedToDelete"] = AuthorizedToDelete(ctx, resource)
	maps["AuthorizedToRestore"] = AuthorizedToRestore(ctx, resource)
	maps["AuthorizedToForceDelete"] = AuthorizedToForceDelete(ctx, resource)

	//// DetailRouterName
	//warp.WithMeta("DetailRouterName", this.DetailRouterName())
	//// EditRouterName
	//warp.WithMeta("EditRouterName", this.EditRouterName())

	isSoftDeleted := false
	if softable, ok := resource.Model().(model.IModel); ok {
		isSoftDeleted = softable.IsSoftDeleted()
	}
	maps["SoftDeleted"] = isSoftDeleted

	// TODO 处理自定义页面返回值
	var item []contracts.Field
	for _, field := range resolveIndexFields(ctx, resource) {
		field.Resolve(ctx, resource.Model())
		item = append(item, field)

		if id, ok := field.(*fields.ID); ok {
			maps["id"] = id
		}
	}

	maps["fields"] = item
	return maps
}

// 详情页数据格式
func SerializeForDetail(ctx *gin.Context, resource contracts.Resource) map[string]interface{} {
	var maps = map[string]interface{}{}
	maps["AuthorizedToUpdate"] = AuthorizedToUpdate(ctx, resource)
	maps["AuthorizedToDelete"] = AuthorizedToDelete(ctx, resource)
	maps["AuthorizedToRestore"] = AuthorizedToRestore(ctx, resource)
	maps["AuthorizedToForceDelete"] = AuthorizedToForceDelete(ctx, resource)

	var isSoftDeleted bool
	if softable, ok := resource.Model().(model.IModel); ok {
		isSoftDeleted = softable.IsSoftDeleted()
	}
	maps["SoftDeleted"] = isSoftDeleted

	var items []contracts.Field
	var p []*panels.Panel
	defaultPanel := panels.NewPanel(resource.Title() + "" + "详情")
	defaultPanel.ShowToolbar = true
	detailFields, panels := resolveDetailFields(ctx, resource)
	data := map[string]interface{}{}
	for _, field := range detailFields {
		if field.GetPanel() == "" {
			defaultPanel.PrepareFields(field)
		}
		field.Resolve(ctx, resource.Model())
		items = append(items, field)

		if id, ok := field.(*fields.ID); ok {
			data["id"] = id
		}
	}

	p = append(p, defaultPanel)
	p = append(p, panels...)

	data["fields"] = items
	maps["panels"] = p
	maps["resource"] = data

	return maps
}

// 资源URI KEY
func ResourceUriKey(resource contracts.Resource) string {
	if customUri, ok := resource.(contracts.ResourceCustomUri); ok {
		return customUri.UriKey()
	} else {
		return utils.StructNameToSnakeAndPlural(resource)
	}
}

// 资源URI PARAM
func ResourceIdParam(resource contracts.Resource) string {
	return utils.StrToSingular(utils.StructToName(resource))
}
