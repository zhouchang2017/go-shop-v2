package vue

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/request"
)

// 列表页api接口
type ResourceHttpIndex interface {
	ResourceHttpIndex() bool
	IndexQuery(ctx *gin.Context, request *request.IndexRequest) error
}

// 详情页api接口
type ResourceHttpShow interface {
	ResourceHttpShow() bool
}

type CustomResourceHttpShow interface {
	CustomResourceHttpShow(ctx *gin.Context, id string) (model interface{}, err error)
}

type ResourceHttpUpdate interface {
	ResourceHttpUpdate() bool
	UpdateFormParse(ctx *gin.Context, model interface{}) (entity interface{}, err error)
}

type ResourceHttpCreate interface {
	ResourceHttpCreate() bool
	CreateFormParse(ctx *gin.Context) (entity interface{}, err error)
}

type ResourceHttpDelete interface {
	ResourceHttpDelete() bool
}

type ResourceHttpForceDelete interface {
	ResourceHttpForceDelete() bool
}

type ResourceHttpRestore interface {
	ResourceHttpRestore() bool
}

// 自定义路由
type CustomHttpRouter interface {
	CustomHttpRouters(router gin.IRouter, uri string, singularLabel string)
}

// api路由生成
func (this *ResourceWarp) httpHandler(router gin.IRouter) {
	this.resourceIndexHandle(router)
	this.resourceDetailHandle(router)
	this.resourceCreateHandle(router)
	this.resourceUpdateHandle(router)
	this.resourceDeleteHandle(router)
	this.resourceForceDeleteHandle(router)
	this.resourceRestoreHandle(router)
	this.resourceCreationFieldHandle(router)
	if custom, ok := this.resource.(CustomHttpRouter); ok {
		custom.CustomHttpRouters(router, this.UriKey(), this.SingularLabel())
	}
}
