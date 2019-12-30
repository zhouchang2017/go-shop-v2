package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"net/http"
	"reflect"
)

type resourceHttpHandle struct {
	ctx      *gin.Context
	router   gin.IRouter
	resource contracts.Resource
	uriKey   string
	idParam  string
}

func (this *resourceHttpHandle) httpHandle() {
	this.resourceIndexHandle()
	this.resourceDetailHandle()
	this.resourceUpdateHandle()
}

func newResourceHttpHandle(router gin.IRouter, resource contracts.Resource) *resourceHttpHandle {
	return &resourceHttpHandle{
		router:   router,
		resource: resource,
		uriKey:   ResourceUriKey(resource),
		idParam:  ResourceIdParam(resource),
	}
}

// vue 路由
func vueRouter() {

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
			filter := &request.IndexRequest{}
			if err := c.ShouldBind(filter); err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			res, err, pagination := paginationable.Pagination(c, filter)

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

// 资源详情api
func (this *resourceHttpHandle) resourceDetailHandle() {
	if showable, ok := this.resource.(contracts.ResourceShowable); ok {
		this.router.GET(fmt.Sprintf("%s/:%s", this.uriKey, this.idParam), func(c *gin.Context) {

			// 验证权限
			if !AuthorizedToView(c, this.resource.Make(nil)) {
				c.AbortWithStatus(403)
				return
			}

			res, err := showable.Show(c, this.idParam)

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			newResource := this.resource.Make(res)

			c.JSON(http.StatusOK, SerializeForDetail(c, newResource))
		})
	}
}

// 资源更新api
func (this *resourceHttpHandle) resourceUpdateHandle() {
	if upgradeable, ok := this.resource.(contracts.ResourceUpgradeable); ok {
		this.router.PUT(fmt.Sprintf("%s/:%s", this.uriKey, this.idParam), func(c *gin.Context) {
			resource := this.resource.Make(nil)
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

			res, err := upgradeable.Update(c, data)

			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			c.JSON(http.StatusOK, gin.H{"redirect": res})
		})
	}
}
