package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/filters"
	"net/http"
	"reflect"
)

func newHttpHandle(vue *Vue, router gin.IRouter) *httpHandle {
	return &httpHandle{vue: vue, router: router}
}

// api路由处理
type httpHandle struct {
	vue    *Vue
	router gin.IRouter
}

func (this *httpHandle) exec() {
	// 授权路由
	this.loginHttpHandle()

	// 授权中间件
	this.router.Use(auth.AuthMiddleware(this.vue.guard))

	// 用户信息路由
	this.userInfoHttpHandle()

	// 系统配置信息
	this.resourcesConfigHttpHandle()

	// vue路由表
	this.vueRoutersHttpHandle()

	// 资源api
	this.resourcesHttpHandle()

	// 自定义页面路由
	this.customPagesHttpHandle()

	// 自定义路由
	this.customHttpHandle()
}

// vue 路由
func (this *httpHandle) vueRoutersHttpHandle() {
	this.router.GET("/routers", func(c *gin.Context) {
		routers := []contracts.Router{}
		for _, warp := range this.vue.warps {
			routers = append(routers, warp.vueRouterFactory.make(c)...)
		}

		// 自定义pages路由
		for _, page := range this.vue.pages {
			if page.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
				routers = append(routers, page.VueRouter())
			}
		}

		c.JSON(http.StatusOK, routers)
	})
}

// api 路由
func (this *httpHandle) resourcesHttpHandle() {
	for _, warp := range this.vue.warps {
		warp.httpHandler.exec(this.router)
	}
}

type adminCredentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// 登录路由
func (this *httpHandle) loginHttpHandle() {
	this.router.POST("/auth/login", func(c *gin.Context) {
		form := &adminCredentials{}
		if err := c.ShouldBind(form); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		jwtGuard, err := auth.Auth.Guard(this.vue.guard)
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
}

// 资源配置信息
func (this *httpHandle) resourcesConfigHttpHandle() {
	this.router.GET("/config", func(c *gin.Context) {
		res := []interface{}{}
		for _, warp := range this.vue.warps {
			info := map[string]interface{}{}
			info["resourceName"] = ResourceUriKey(warp.resource)
			info["title"] = warp.resource.Title()
			res = append(res, info)
		}
		c.JSON(http.StatusOK, res)
	})
}

// 用户信息路由
func (this *httpHandle) userInfoHttpHandle() {
	this.router.GET("/auth/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, ctx.GetUser(c))
	})
}

// 自定义页面路由
func (this *httpHandle) customPagesHttpHandle() {
	for _, page := range this.vue.pages {
		page.HttpHandles(this.router)
	}
}

// 自定义路由
func (this *httpHandle) customHttpHandle() {
	for _, handle := range this.vue.customHttpHandle {
		handle(this.router)
	}
}
func newResourceHttpHandle(resource contracts.Resource) *resourceHttpHandle {
	return &resourceHttpHandle{
		resource:      resource,
		uriKey:        ResourceUriKey(resource),
		idParam:       ResourceIdParam(resource),
		isSoftDeleted: ResourceIsSoftDeleted(resource),
	}
}

type resourceHttpHandle struct {
	router        gin.IRouter
	resource      contracts.Resource
	uriKey        string
	idParam       string
	isSoftDeleted bool
}

func (this *resourceHttpHandle) exec(router gin.IRouter) {
	this.router = router
	this.resourceIndexHandle()         // 列表
	this.resourceIndexApiHandle()      // 列表api
	this.resourceDetailHandle()        // 详情
	this.resourceDetailApiHandle()     // 详情api
	this.resourceUpdateHandle()        // 更新
	this.resourceCreateHandle()        // 创建
	this.resourceCreationFieldHandle() // 创建字段
	this.resourceUpdateFieldHandle()   // 更新字段
	this.resourceDestroyHandle()       // 删除
	this.resourceForceDestroyHandle()  // 销毁
	this.resourceRestoreHandle()       // 还原

	this.resourceLensesHandle() // 聚合

	this.resourcePagesHandle() // 自定义页面

	this.resourceFiltersHandle() // 过滤
}

// 资源列表页api
func (this *resourceHttpHandle) resourceIndexHandle() {
	if paginationable, ok := this.resource.(contracts.ResourcePaginationable); ok {

		this.router.GET(this.uriKey, func(c *gin.Context) {

			// 验证权限
			if !AuthorizedToViewAny(c, this.resource.Make(nil)) {
				c.AbortWithStatus(403)
				return
			}

			// 处理函数
			form := &request.IndexRequest{}
			if err := c.ShouldBind(form); err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			filters2 := form.Filters.Unmarshal()

			// filters处理
			for _, filter := range this.resource.Filters(c) {
				if err := filter.Apply(c, filters2[filter.Key()], form); err != nil {
					err2.ErrorEncoder(nil, err, c.Writer)
					return
				}
			}

			res, pagination, err := paginationable.Pagination(c, form)

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			indexResources := []interface{}{}
			if reflect.TypeOf(res).Kind() == reflect.Slice {
				valueOf := reflect.ValueOf(res)
				len := valueOf.Len()
				for i := 0; i < len; i++ {
					model := valueOf.Index(i).Interface()
					resource := this.resource.Make(model)
					indexResources = append(indexResources, SerializeForIndex(c, resource))
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"pagination": pagination,
				"data":       indexResources,
			})
		})
	}
}

func (this *resourceHttpHandle) resourceIndexApiHandle() {
	if paginationable, ok := this.resource.(contracts.ResourcePaginationable); ok {

		this.router.GET(fmt.Sprintf("api/%s", this.uriKey), func(c *gin.Context) {

			// 验证权限
			if !AuthorizedToViewAny(c, this.resource.Make(nil)) {
				c.AbortWithStatus(403)
				return
			}

			// 处理函数
			form := &request.IndexRequest{}
			if err := c.ShouldBind(form); err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			filters2 := form.Filters.Unmarshal()

			// filters处理
			for _, filter := range this.resource.Filters(c) {
				if err := filter.Apply(c, filters2[filter.Key()], form); err != nil {
					err2.ErrorEncoder(nil, err, c.Writer)
					return
				}
			}

			res, pagination, err := paginationable.Pagination(c, form)

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"pagination": pagination,
				"data":       res,
			})
		})
	}
}

// 资源详情后台接口
func (this *resourceHttpHandle) resourceDetailHandle() {
	if showable, ok := this.resource.(contracts.ResourceShowable); ok {
		this.router.GET(fmt.Sprintf("%s/:%s", this.uriKey, this.idParam), func(c *gin.Context) {

			// 验证权限
			if !AuthorizedToView(c, this.resource.Make(nil)) {
				c.AbortWithStatus(403)
				return
			}

			res, err := showable.Show(c, c.Param(this.idParam))

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			newResource := this.resource.Make(res)

			c.JSON(http.StatusOK, SerializeForDetail(c, newResource))
		})
	}
}

// 资源详情api
func (this *resourceHttpHandle) resourceDetailApiHandle() {
	if showable, ok := this.resource.(contracts.ResourceShowable); ok {
		this.router.GET(fmt.Sprintf("api/%s/:%s", this.uriKey, this.idParam), func(c *gin.Context) {

			// 验证权限
			if !AuthorizedToView(c, this.resource.Make(nil)) {
				c.AbortWithStatus(403)
				return
			}

			res, err := showable.Show(c, c.Param(this.idParam))

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			c.JSON(http.StatusOK, res)
		})
	}
}

// 资源更新后台接口
func (this *resourceHttpHandle) resourceUpdateHandle() {
	if showable, ok := this.resource.(contracts.ResourceShowable); ok {
		if upgradeable, ok := this.resource.(contracts.ResourceUpgradeable); ok {
			this.router.PUT(fmt.Sprintf("%s/:%s", this.uriKey, this.idParam), func(c *gin.Context) {
				model, err := showable.Show(c, c.Param(this.idParam))
				if err != nil {
					// 404
					err2.ErrorEncoder(nil, err, c.Writer)
					return
				}

				resource := this.resource.Make(model)

				// 验证权限
				if !AuthorizedToUpdate(c, resource) {
					c.AbortWithStatus(403)
					return
				}

				// 表单验证
				fields, _ := resolveUpdateFields(c, resource)
				data, err := Validator(c, fields)
				if err != nil {
					err2.ErrorEncoder(nil, err, c.Writer)
					return
				}

				res, err := upgradeable.Update(c, model, data)

				if err != nil {
					err2.ErrorEncoder(nil, err, c.Writer)
					return
				}

				c.JSON(http.StatusOK, gin.H{"redirect": res})
			})
		}
	}

}

// 更新资源字段
func (this *resourceHttpHandle) resourceUpdateFieldHandle() {
	if showable, ok := this.resource.(contracts.ResourceShowable); ok {
		if _, ok := this.resource.(contracts.ResourceUpgradeable); ok {
			this.router.GET(fmt.Sprintf("update-fields/%s/:%s", this.uriKey, this.idParam), func(c *gin.Context) {
				res, err := showable.Show(c, c.Param(this.idParam))
				if err != nil {
					// 404
					err2.ErrorEncoder(nil, err, c.Writer)
					return
				}

				resource := this.resource.Make(res)
				// 验证权限
				if !AuthorizedToUpdate(c, resource) {
					c.AbortWithStatus(403)
					return
				}

				fields, panels := resolveUpdateFields(c, resource)

				for _, field := range fields {
					field.Call()

					field.Resolve(c, resource.Model())
				}

				c.JSON(http.StatusOK, gin.H{
					"fields": fields,
					"panels": panels,
				})
			})
		}
	}

}

// 资源创建api
func (this *resourceHttpHandle) resourceCreateHandle() {
	// 如果实现自定义创建页
	if _, ok := this.resource.(contracts.ResourceCustomCreationComponent); ok {
		return
	}
	if storeable, ok := this.resource.(contracts.ResourceStorable); ok {
		this.router.POST(fmt.Sprintf("%s", this.uriKey), func(c *gin.Context) {
			resource := this.resource.Make(nil)
			// 验证权限
			if !AuthorizedToCreate(c, resource) {
				c.AbortWithStatus(403)
				return
			}

			// 表单验证
			fields, _ := resolveCreationFields(c, resource)
			data, err := Validator(c, fields)
			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			res, err := storeable.Store(c, data)

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			c.JSON(http.StatusCreated, gin.H{"redirect": res})
		})
	}
}

// 创建资源字段
func (this *resourceHttpHandle) resourceCreationFieldHandle() {
	// 如果实现自定义创建页
	if _, ok := this.resource.(contracts.ResourceCustomCreationComponent); ok {
		return
	}

	if _, ok := this.resource.(contracts.ResourceStorable); ok {
		this.router.GET(fmt.Sprintf("creation-fields/%s", this.uriKey), func(c *gin.Context) {
			resource := this.resource.Make(nil)
			// 验证权限
			if !AuthorizedToCreate(c, resource) {
				c.AbortWithStatus(403)
				return
			}

			fields, panels := resolveCreationFields(c, resource)

			for _, field := range fields {
				field.Call()
			}

			c.JSON(http.StatusOK, gin.H{
				"fields": fields,
				"panels": panels,
			})
		})
	}
}

// 资源删除api
func (this *resourceHttpHandle) resourceDestroyHandle() {
	if destroyable, ok := this.resource.(contracts.ResourceDestroyable); ok {
		this.router.DELETE(fmt.Sprintf("%s/:%s", this.uriKey, this.idParam), func(c *gin.Context) {
			resource := this.resource.Make(nil)
			// 验证权限
			if !AuthorizedToDelete(c, resource) {
				c.AbortWithStatus(403)
				return
			}

			err := destroyable.Destroy(c, c.Param(this.idParam))

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			c.JSON(http.StatusOK, nil)
		})
	}
}

// 资源销毁api
func (this *resourceHttpHandle) resourceForceDestroyHandle() {
	if destroyable, ok := this.resource.(contracts.ResourceForceDestroyable); ok && this.isSoftDeleted {
		this.router.DELETE(fmt.Sprintf("%s/:%s/force", this.uriKey, this.idParam), func(c *gin.Context) {
			resource := this.resource.Make(nil)
			// 验证权限
			if !AuthorizedToDelete(c, resource) {
				c.AbortWithStatus(403)
				return
			}

			err := destroyable.ForceDestroy(c, c.Param(this.idParam))

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			c.JSON(http.StatusOK, nil)
		})
	}
}

// 还原资源api
func (this *resourceHttpHandle) resourceRestoreHandle() {
	if restoreable, ok := this.resource.(contracts.ResourceRestoreable); ok && this.isSoftDeleted {
		this.router.DELETE(fmt.Sprintf("%s/:%s/force", this.uriKey, this.idParam), func(c *gin.Context) {
			resource := this.resource.Make(nil)
			// 验证权限
			if !AuthorizedToRestore(c, resource) {
				c.AbortWithStatus(403)
				return
			}

			err := restoreable.Restore(c, c.Param(this.idParam))

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			c.JSON(http.StatusOK, nil)
		})
	}
}

// lens聚合api
func (this *resourceHttpHandle) resourceLensesHandle() {

	// 获取所有聚合
	this.router.GET(fmt.Sprintf("/lenses/%s", this.uriKey), func(c *gin.Context) {
		data := []interface{}{}
		for _, lens := range this.resource.Lenses() {
			// 验证权限
			item := map[string]interface{}{}
			if lens.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
				item["router_name"] = LensRouteName(this.resource, lens)
				item["title"] = lens.Title()
			}
			data = append(data, item)
		}
		c.JSON(http.StatusOK, data)
	})

	for _, lens := range this.resource.Lenses() {
		this.router.GET(LensEndPoints(this.resource, lens), func(c *gin.Context) {
			// 验证权限
			if !lens.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
				c.AbortWithStatus(403)
				return
			}

			// 处理函数
			form := &request.IndexRequest{}
			if err := c.ShouldBind(form); err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}
			filters2 := form.Filters.Unmarshal()
			// filters处理
			for _, filter := range lens.Filters(c) {
				if err := filter.Apply(c, filters2[filter.Key()], form); err != nil {
					err2.ErrorEncoder(nil, err, c.Writer)
					return
				}
			}

			res, pagination, err := lens.HttpHandle(c, form)

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			data := []interface{}{}
			if reflect.TypeOf(res).Kind() == reflect.Slice {
				valueOf := reflect.ValueOf(res)
				len := valueOf.Len()
				for i := 0; i < len; i++ {
					model := valueOf.Index(i).Interface()
					data = append(data, SerializeForLensIndex(c, lens, model))
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"pagination": pagination,
				"data":       data,
			})

		})

		// 获取聚合过滤
		this.router.GET(LensFiltersEndPoints(this.resource, lens), func(c *gin.Context) {
			c.JSON(http.StatusOK, serializeFiltersForMaps(c, resolverFilters(lens, c)...))
		})
	}
}

// 自定义页面
func (this *resourceHttpHandle) resourcePagesHandle() {
	// 获取所有自定义页面
	this.router.GET(fmt.Sprintf("/pages/%s", this.uriKey), func(c *gin.Context) {
		data := []interface{}{}
		for _, page := range this.resource.Pages() {
			// 验证权限
			item := map[string]interface{}{}
			if page.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
				item["router_name"] = page.VueRouter().RouterName()
				item["title"] = page.Title()
			}
			data = append(data, item)
		}
		c.JSON(http.StatusOK, data)
	})

}

// 过滤选项
func (this *resourceHttpHandle) resourceFiltersHandle() {
	this.router.GET(fmt.Sprintf("/filters/%s", this.uriKey), func(c *gin.Context) {
		c.JSON(http.StatusOK, serializeFiltersForMaps(c, resolverFilters(this.resource, c)...))
	})
}

func serializeFiltersForMaps(ctx *gin.Context, f ...contracts.Filter) []map[string]interface{} {
	data := []map[string]interface{}{}
	for _, filter := range f {
		data = append(data, filters.SerializeMap(ctx, filter))
	}
	return data
}

func resolverFilters(target interface{ Filters(ctx *gin.Context) []contracts.Filter }, c *gin.Context) []contracts.Filter {
	filters := []contracts.Filter{}
	for _, filter := range target.Filters(c) {
		// 验证权限
		if filter.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			filters = append(filters, filter)
		}
	}
	return filters
}
