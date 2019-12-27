package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/repository"
)

// 可展示icon图标
type Iconable interface {
	Icon() string
}

type Resource interface {
	Group() string                             // 侧边栏分组
	DisplayInNavigation(ctx *gin.Context) bool // 导航栏是否显示
	Model() interface{}                        // 实例对象
	Repository() repository.IRepository        // 数据仓库
	Make(model interface{}) Resource           // 实例化方法
	SetModel(model interface{})
	SetRoot(vue *Vue)
	Title() string
	// 前端路由定义
	HasIndexRoute(ctx *gin.Context) bool
	HasDetailRoute(ctx *gin.Context) bool
	HasEditRoute(ctx *gin.Context) bool

	Lenses() []Lens // 自定义聚合查询等
	Links() []Link  // 自定义Link
	Observer
}

type ResourceWarp struct {
	resource Resource
	root     *Vue
	*ResourceHelper
}

func NewResourceWarp(resource Resource, root *Vue) *ResourceWarp {
	return &ResourceWarp{resource: resource, root: root, ResourceHelper: NewResourceHelper(resource)}
}

func (this *ResourceWarp) CreateButtonName() string {
	if custom, ok := this.resource.(CustomCreateButtonName); ok {
		return custom.CreateButtonName()
	}
	return fmt.Sprintf("创建%s", this.resource.Title())
}

// 生成前端路由对象
func (this *ResourceWarp) routers(ctx *gin.Context) []*Router {
	var routers []*Router
	uri := this.VueUriKey()

	var authorizedToCreate = this.AuthorizedToCreate(ctx)
	// 自定义路由
	if customVueRouter, ok := this.resource.(CustomVueRouter); ok {
		router := customVueRouter.CustomVueRouter(ctx, this)
		routers = append(routers, router...)
	}
	// 列表页路由
	if this.resource.HasIndexRoute(ctx) {
		router := &Router{
			Path: uri,
			Name: this.IndexRouterName(),
			//Component: fmt.Sprintf(`%s/Index`, uri),
			Component: "Index",
			Hidden:    !this.resource.DisplayInNavigation(ctx),
		}

		router.WithMeta("AuthorizedToCreate", authorizedToCreate)
		router.WithMeta("Title", this.resource.Title())
		router.WithMeta("ResourceName", this.UriKey())

		router.WithMeta("CreateButtonText", this.CreateButtonName())
		router.WithMeta("CreateRouterName", this.CreateRouterName())
		router.WithMeta("DetailRouterName", this.DetailRouterName())
		router.WithMeta("EditRouterName", this.EditRouterName())
		router.WithMeta("Group", this.resource.Group())
		if iconable, ok := this.resource.(Iconable); ok {
			router.WithMeta("icon", iconable.Icon())
		}
		// 追加列
		router.WithMeta("Headings", this.resolveIndexFields(ctx))

		if listener, ok := this.resource.(ListenerIndexRouteCreated); ok {
			listener.OnIndexRouteCreated(ctx, router)
		}
		routers = append(routers, router)
	}

	// 创建页面路由
	if authorizedToCreate {
		router := &Router{
			Path:      fmt.Sprintf("%s/create", uri),
			Name:      this.CreateRouterName(),
			Component: fmt.Sprintf(`%s/Create`, uri),
			Hidden:    true,
		}
		router.WithMeta("Title", this.CreateButtonName())
		router.WithMeta("DetailRouterName", this.DetailRouterName())
		router.WithMeta("IndexRouterName", this.IndexRouterName())
		if listener, ok := this.resource.(ListenerCreateRouteCreated); ok {
			listener.OnCreateRouteCreated(ctx, router)
		}
		routers = append(routers, router)
	}

	// 详情页路由
	if this.resource.HasDetailRoute(ctx) {
		router := &Router{
			Path: fmt.Sprintf("%s/:id", uri),
			Name: this.DetailRouterName(),
			//Component: fmt.Sprintf(`%s/Detail`, uri),
			Component: "Detail",
			Hidden:    true,
		}
		router.WithMeta("Title", this.resource.Title()+""+"详情")
		router.WithMeta("ResourceName", this.UriKey())
		router.WithMeta("IndexRouterName", this.IndexRouterName())
		if listener, ok := this.resource.(ListenerDetailRouteCreated); ok {
			listener.OnDetailRouteCreated(ctx, router)
		}
		routers = append(routers, router)
	}

	// 更新页面路由
	if this.resource.HasEditRoute(ctx) {
		router := &Router{
			Path:      fmt.Sprintf("%s/:id/edit", uri),
			Name:      this.EditRouterName(),
			Component: fmt.Sprintf(`%s/Edit`, uri),
			Hidden:    true,
		}
		router.WithMeta("Title", this.resource.Title()+""+"编辑")
		router.WithMeta("DetailRouterName", this.DetailRouterName())
		router.WithMeta("IndexRouterName", this.IndexRouterName())
		if listener, ok := this.resource.(ListenerUpdateRouteCreated); ok {
			listener.OnUpdateRouteCreated(ctx, router)
		}
		routers = append(routers, router)
	}

	// lenses routers
	for _, lens := range this.resource.Lenses() {
		if lens.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			router := &Router{
				Path:      fmt.Sprintf("%s/lenses/%s", uri, lens.RouterName()),
				Name:      lensRouterName(lens, uri),
				Component: lens.Component(),
				Hidden:    true,
			}
			router.WithMeta("Title", lens.Title())
			router.WithMeta("IndexRouterName", this.IndexRouterName())
			router.WithMeta("ResourceName", this.SingularLabel())
			router.WithMeta("IndexTitle", this.resource.Title())
			router.WithMeta("LensApiUri", lensApiUri(lens, this.UriKey()))

			if listener, ok := lens.(ListenerLensRouteCreated); ok {
				listener.OnLensRouteCreated(ctx, router)
			}

			routers = append(routers, router)
		}
	}

	return routers
}

// 列表页数据格式
func (this *ResourceWarp) SerializeForIndex(ctx *gin.Context) Metable {
	warp := &responseWarp{}
	var maps = map[string]bool{}
	maps["AuthorizedToView"], _ = this.AuthorizedToView(ctx)
	maps["AuthorizedToUpdate"], _ = this.AuthorizedToUpdate(ctx)
	maps["AuthorizedToDelete"], _ = this.AuthorizedToDelete(ctx)
	maps["AuthorizedToRestore"], _ = this.AuthorizedToRestore(ctx)
	maps["AuthorizedToForceDelete"], _ = this.AuthorizedToForceDelete(ctx)

	for k, v := range maps {
		warp.WithMeta(k, v)
	}
	// DetailRouterName
	warp.WithMeta("DetailRouterName", this.DetailRouterName())
	// EditRouterName
	warp.WithMeta("EditRouterName", this.EditRouterName())

	isSoftDeleted := false
	if softable, ok := this.resource.Model().(model.IModel); ok {
		isSoftDeleted = softable.IsSoftDeleted()
	}
	warp.WithMeta("SoftDeleted", isSoftDeleted)

	if _, ok := this.resource.(HasFields); ok {

		var item []Field
		for _, field := range this.resolveIndexFields(ctx) {
			field.Resolve(ctx, this.resource.Model())
			item = append(item, field)

			if id, ok := field.(*ID); ok {
				warp.Id = id
			}
		}

		warp.Data = item

	} else {
		warp.Data = this.resource.Model()
	}

	return warp
}

func (this *ResourceWarp) resolveIndexFields(ctx *gin.Context) []Field {
	var item []Field
	if hasFields, ok := this.resource.(HasFields); ok {
		fields := hasFields.Fields(ctx, this.resource.Model())
		for _, field := range fields() {
			if isField, ok := field.(Field); ok {
				if isField.ShowOnIndex() && isField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
					item = append(item, isField)
					continue
				}
			}

			if isPanel, ok := field.(*Panel); ok {
				for _, panelField := range isPanel.Fields {
					if panelField.ShowOnIndex() && panelField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
						item = append(item, panelField)
					}
				}
			}

		}
	}
	return item
}

func (this *ResourceWarp) resolveDetailFields(ctx *gin.Context) ([]Field, []*Panel) {
	var item []Field
	var panel []*Panel
	if hasFields, ok := this.resource.(HasFields); ok {
		fields := hasFields.Fields(ctx, this.resource.Model())
		for _, field := range fields() {

			if isField, ok := field.(Field); ok {
				if isField.ShowOnDetail() && isField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
					item = append(item, isField)
					continue
				}
			}

			if isPanel, ok := field.(*Panel); ok {
				for _, panelField := range isPanel.Fields {
					if panelField.ShowOnDetail() && panelField.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
						item = append(item, panelField)
					}
				}
				panel = append(panel, isPanel)
			}

		}
	}
	return item, panel
}

type responseWarp struct {
	Meta MetaItems   `json:"meta"`
	Id   *ID         `json:"id"`
	Data interface{} `json:"fields"`
}

func (m *responseWarp) WithMeta(key string, value interface{}) {
	m.Meta = append(m.Meta, &metaItem{key, value})
}

type detailResourceWarp struct {
	Meta     MetaItems `json:"meta"`
	Panels   []*Panel  `json:"panels"`
	Resource struct {
		Id     *ID     `json:"id"`
		Fields []Field `json:"fields"`
	} `json:"resource"`
	Data interface{} `json:"data,omitempty"`
}

func (m *detailResourceWarp) WithMeta(key string, value interface{}) {
	m.Meta = append(m.Meta, &metaItem{key, value})
}

// 详情页数据格式
func (this *ResourceWarp) SerializeForDetail(ctx *gin.Context) Metable {
	warp := &detailResourceWarp{}
	var maps = map[string]bool{}
	maps["AuthorizedToUpdate"], _ = this.AuthorizedToUpdate(ctx)
	maps["AuthorizedToDelete"], _ = this.AuthorizedToDelete(ctx)
	maps["AuthorizedToRestore"], _ = this.AuthorizedToRestore(ctx)
	maps["AuthorizedToForceDelete"], _ = this.AuthorizedToForceDelete(ctx)

	for k, v := range maps {
		warp.WithMeta(k, v)
	}

	// with edit router
	if (maps["AuthorizedToUpdate"]) {
		warp.WithMeta("EditRouterName", this.EditRouterName())
	}

	isSoftDeleted := false
	if softable, ok := this.resource.Model().(model.IModel); ok {
		isSoftDeleted = softable.IsSoftDeleted()
	}
	warp.WithMeta("SoftDeleted", isSoftDeleted)

	if _, ok := this.resource.(HasFields); ok {
		var items []Field
		var p []*Panel
		defaultPanel := NewPanel(this.resource.Title() + "" + "详情")
		defaultPanel.ShowToolbar = true
		fields, panels := this.resolveDetailFields(ctx)
		for _, field := range fields {
			if field.GetPanel() == "" {
				defaultPanel.PrepareFields(field)
			}
			field.Resolve(ctx, this.resource.Model())
			items = append(items, field)

			if id, ok := field.(*ID); ok {
				warp.Resource.Id = id
			}
		}

		p = append(p, defaultPanel)
		p = append(p, panels...)

		warp.Resource.Fields = items
		warp.Panels = p
	} else {
		warp.Data = this.resource.Model()
	}

	return warp
}
