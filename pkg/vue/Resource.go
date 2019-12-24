package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
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
			Path:      uri,
			Name:      this.IndexRouterName(),
			Component: fmt.Sprintf(`%s/Index`, uri),
			Hidden:    !this.resource.DisplayInNavigation(ctx),
		}

		router.WithMeta("AuthorizedToCreate", authorizedToCreate)
		router.WithMeta("Title", this.resource.Title())
		router.WithMeta("ResourceName", this.SingularLabel())
		router.WithMeta("ResourceUriKey", this.UriKey())

		router.WithMeta("CreateButtonText", this.CreateButtonName())
		router.WithMeta("CreateRouterName", this.CreateRouterName())
		router.WithMeta("DetailRouterName", this.DetailRouterName())
		router.WithMeta("EditRouterName", this.EditRouterName())
		router.WithMeta("Group", this.resource.Group())
		if iconable, ok := this.resource.(Iconable); ok {
			router.WithMeta("icon", iconable.Icon())
		}

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
			Path:      fmt.Sprintf("%s/:id", uri),
			Name:      this.DetailRouterName(),
			Component: fmt.Sprintf(`%s/Detail`, uri),
			Hidden:    true,
		}
		router.WithMeta("Title", this.resource.Title()+""+"详情")
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
	warp.Data = this.resource.Model()
	return warp
}

type responseWarp struct {
	Meta MetaItems   `json:"meta"`
	Data interface{} `json:"data"`
}

func (m *responseWarp) WithMeta(key string, value interface{}) {
	m.Meta = append(m.Meta, &metaItem{key, value})
}

type detailResourceWarp struct {
	AuthorizedToUpdate      bool
	AuthorizedToDelete      bool
	AuthorizedToRestore     bool
	AuthorizedToForceDelete bool
	Data                    interface{} `json:"data"`
}

// 详情页数据格式
func (this *ResourceWarp) SerializeForDetail(ctx *gin.Context) Metable {
	warp := &responseWarp{}
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

	warp.Data = this.resource.Model()
	return warp
}
