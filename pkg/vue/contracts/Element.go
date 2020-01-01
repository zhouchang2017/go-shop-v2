package contracts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
)

type Metable interface {
	WithMeta(key string, value interface{})
}

// 前端页面元素
type Element interface {
	Component() string

	PrefixComponent() bool
	// 是否有权限可见
	AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool
	// 附加值
	WithMeta(key string, value interface{})
	// meta
	Meta() map[string]interface{}
}
