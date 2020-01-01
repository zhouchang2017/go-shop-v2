package contracts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type Lens interface {
	Title() string                                                 // 标题
	HttpHandle(ctx *gin.Context, request *request.IndexRequest) (res interface{}, pagination response.Pagination, err error)    // http处理
	AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool // 授权
	// 字段
	Fields(ctx *gin.Context,model interface{}) func() []interface{}
	// 过滤
	Filters(ctx *gin.Context) []Filter
}

