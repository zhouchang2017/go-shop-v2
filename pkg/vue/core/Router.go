package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/vue/contracts"
)

type Router struct {
	RouterPath      string                 `json:"path"`
	RouterComponent string                 `json:"component"`
	Name            string                 `json:"name,omitempty"` // 命名路由
	Children        []contracts.Router     `json:"children,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
	Hidden          bool                   `json:"hidden"`
}

func NewRouter() *Router {
	return &Router{Meta: map[string]interface{}{}}
}

func (m Router) Path() string {
	return m.RouterPath
}

func (m Router) Component() string {
	return m.RouterComponent
}

func (m Router) RouterName() string {
	return m.Name
}

func (m *Router) WithMeta(key string, value interface{}) {
	m.Meta[key] = value
}

func (m *Router) AddChild(r contracts.Router) {
	m.Children = append(m.Children, r)
}

// 列表页路由名称
func IndexRouteName(resource contracts.Resource) string {
	return fmt.Sprintf("%s.index", ResourceUriKey(resource))
}

// 详情页路由名称
func DetailRouteName(resource contracts.Resource) string {
	return fmt.Sprintf("%s.detail", ResourceUriKey(resource))
}

// 更新页路由名称
func UpdateRouteName(resource contracts.Resource) string {
	return fmt.Sprintf("%s.edit", ResourceUriKey(resource))
}

// 创建页路由名称
func CreationRouteName(resource contracts.Resource) string {
	return fmt.Sprintf("%s.create", ResourceUriKey(resource))
}

// 列表页组件
func IndexRouteComponent(resource contracts.Resource) string {
	if implement, ok := resource.(contracts.ResourceCustomIndexComponent); ok {
		return implement.IndexComponent().VueRouter().Component()
	}
	return "Index"
}

// 创建页组件
func CreationRouteComponent(resource contracts.Resource) string {
	if implement, ok := resource.(contracts.ResourceCustomCreationComponent); ok {
		return implement.CreationComponent().VueRouter().Component()
	}
	return "Create"
}

// 详情页组件
func DetailRouteComponent(resource contracts.Resource) string {
	if implement, ok := resource.(contracts.ResourceCustomDetailComponent); ok {
		return implement.DetailComponent().VueRouter().Component()
	}
	return "Detail"
}

// 更新页组件
func UpdateRouteComponent(resource contracts.Resource) string {
	if implement, ok := resource.(contracts.ResourceCustomUpdateComponent); ok {
		return implement.UpdateComponent().VueRouter().Component()
	}
	return "Update"
}

type vueRouterFactory struct {
	ctx                *gin.Context
	user               interface{}
	resource           contracts.Resource
	uriKey             string
	resourceName       string
	indexComponent     string
	detailComponent    string
	updateComponent    string
	createComponent    string
	indexRouteName     string
	detailRouteName    string
	createRouteName    string
	updateRouteName    string
	authorizedToCreate bool
}

func newVueRouterFactory(resource contracts.Resource) *vueRouterFactory {
	return &vueRouterFactory{
		resource:     resource,
		uriKey:       ResourceUriKey(resource),
		resourceName: ResourceUriKey(resource),
		// 组件
		indexComponent:  IndexRouteComponent(resource),
		detailComponent: DetailRouteComponent(resource),
		updateComponent: UpdateRouteComponent(resource),
		createComponent: CreationRouteComponent(resource),
		// 名称
		indexRouteName:  IndexRouteName(resource),
		detailRouteName: DetailRouteName(resource),
		createRouteName: CreationRouteName(resource),
		updateRouteName: UpdateRouteName(resource),
	}
}

func (this *vueRouterFactory) make(ctx *gin.Context) []contracts.Router {
	this.ctx = ctx
	this.user = ctx2.GetUser(ctx)
	this.authorizedToCreate = AuthorizedToCreate(ctx, this.resource)
	var routers []contracts.Router
	if router := this.vueIndexRouter(); router != nil {
		routers = append(routers, router)
	}
	if router := this.vueDetailRouter(); router != nil {
		routers = append(routers, router)
	}
	if router := this.vueUpdateRouter(); router != nil {
		routers = append(routers, router)
	}
	if router := this.vueCreateRouter(); router != nil {
		routers = append(routers, router)
	}
	return routers
}

// vue 资源列表页路由
func (this *vueRouterFactory) vueIndexRouter() contracts.Router {
	if this.resource.HasIndexRoute(this.ctx, this.user) && AuthorizedToViewAny(this.ctx, this.resource) {
		router := NewRouter()
		router.RouterPath = this.uriKey
		router.Name = this.indexRouteName
		router.RouterComponent = this.indexComponent
		router.Hidden = !this.resource.DisplayInNavigation(this.ctx, this.user)

		router.WithMeta("AuthorizedToCreate", this.authorizedToCreate)
		router.WithMeta("Title", this.resource.Title())
		router.WithMeta("ResourceName", this.resourceName)

		router.WithMeta("CreateButtonText", fmt.Sprintf("创建%s", this.resource.Title()))
		router.WithMeta("CreateRouterName", this.createRouteName)
		router.WithMeta("DetailRouterName", this.detailRouteName)
		router.WithMeta("EditRouterName", this.updateRouteName)
		router.WithMeta("Group", this.resource.Group())
		if iconable, ok := this.resource.(contracts.Iconable); ok {
			router.WithMeta("icon", iconable.Icon())
		}
		// 追加列
		router.WithMeta("Headings", resolveIndexFields(this.ctx, this.resource))
		return router
	}
	return nil
}

// vue 资源详情页路由
func (this *vueRouterFactory) vueDetailRouter() contracts.Router {
	if this.resource.HasDetailRoute(this.ctx, this.user) {
		router := NewRouter()
		router.RouterPath = fmt.Sprintf("%s/:id", this.uriKey)
		router.Name = this.detailRouteName
		router.RouterComponent = this.detailComponent
		router.Hidden = true
		router.WithMeta("ResourceName", this.resourceName)
		return router
	}
	return nil
}

// vue 资源更新页路由
func (this *vueRouterFactory) vueUpdateRouter() contracts.Router {

	if _, ok := this.resource.(contracts.ResourceCustomUpdateComponent); ok {
		return nil
	}

	if this.resource.HasEditRoute(this.ctx, this.user) {
		router := NewRouter()
		router.RouterPath = fmt.Sprintf("%s/:id/edit", this.uriKey)
		router.Name = this.updateRouteName
		router.RouterComponent = this.updateComponent
		router.Hidden = true
		router.WithMeta("ResourceName", this.resourceName)
		return router
	}
	return nil
}

// vue 资源创建页路由
func (this *vueRouterFactory) vueCreateRouter() contracts.Router {
	if this.authorizedToCreate {

		if _, ok := this.resource.(contracts.ResourceCustomCreationComponent); ok {
			return nil
		}

		router := NewRouter()
		router.RouterPath = fmt.Sprintf("%s/new", this.uriKey)
		router.Name = this.createRouteName
		router.RouterComponent = this.createComponent
		router.Hidden = true
		router.WithMeta("ResourceName", this.resourceName)
		return router
	}
	return nil
}