package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/utils"
	"net/http"
)

type Vue struct {
	resources        []Resource
	guard            string
	policies         map[string]interface{}
	customHttpHandle []func(router gin.IRouter)
	warps            map[string]*ResourceWarp
}

func (this *Vue) SetGuard(guard string) {
	this.guard = guard
}

func (this *Vue) RegisterCustomHttpHandler(handler func(router gin.IRouter)) {
	this.customHttpHandle = append(this.customHttpHandle, handler)
}

// 注册资源
func (this *Vue) RegisterResource(resource Resource) {
	this.resources = append(this.resources, resource)
}

// 注册策略
func (this *Vue) RegisterPolice(policy interface{}) {
	this.policies[utils.StructToName(policy)] = policy
}

func (this *Vue) resolvePolicy(resource interface{}) (interface{}, bool) {
	if policy, ok := this.policies[utils.StructToName(resource)+"Policy"]; ok {
		return policy, true
	}
	return nil, false
}

type adminCredentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (this *Vue) authHttpHandler(router gin.IRouter) {
	router.POST("/auth/login", func(c *gin.Context) {
		form := &adminCredentials{}
		if err := c.ShouldBind(form); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		jwtGuard, err := auth.Auth.Guard(this.guard)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		credentials := map[string]string{
			"username": form.Username,
			"password": form.Password,
		}
		if res, ok := jwtGuard.Attempt(credentials, true); ok {
			c.Header("token", fmt.Sprintf("%s", res))
			c.JSON(http.StatusOK, res)
			return
		}
		err2.ErrorEncoder(nil,
			err2.NewFromCode(422).F("用户名或密码错误"),
			c.Writer)
	})

	authGroup := router.Use(auth.AuthMiddleware(this.guard))
	authGroup.GET("/auth/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, ctx.GetUser(c))
	})
}

func (this *Vue) addWarp(resource Resource, warp *ResourceWarp) {
	this.warps[utils.StructToName(resource)] = warp
}

func (this *Vue) ResolveWarp(resource Resource) (warp *ResourceWarp, ok bool) {
	if warp, ok = this.warps[utils.StructToName(resource)]; ok {
		return warp, true
	}
	return nil, false
}

func (this *Vue) httpHandler(router gin.IRouter) {

	this.authHttpHandler(router)
	router.Use(auth.AuthMiddleware(this.guard))

	// RESTFUL API
	for _, resource := range this.resources {
		warp := newResourceWarp(resource, this)
		resource.SetRoot(this)

		this.addWarp(resource, warp) // 保存warp

		warp.httpHandler(router)
		warp.resourceLensesIndexHandle(router)
		warp.resourceLensesDetailHandle(router)
	}

	// vue路由配置
	this.providerVueRouteConfig(router)

	// 自定义路由
	for _, handle := range this.customHttpHandle {
		handle(router)
	}
}

func (this *Vue) providerVueRouteConfig(router gin.IRouter) {
	router.GET("/routers", func(c *gin.Context) {
		routers := this.vueRouters(c)
		c.JSON(http.StatusOK, routers)
	})
}

func (this *Vue) vueRouters(ctx *gin.Context) []*Router {
	var routers []*Router
	for _, warp := range this.warps {
		routers = append(routers, warp.routers(ctx)...)
	}
	return routers
}

func NewVue() *Vue {
	return &Vue{
		policies: map[string]interface{}{},
		warps:    map[string]*ResourceWarp{},
	}
}
