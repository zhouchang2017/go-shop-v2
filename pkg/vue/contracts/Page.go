package contracts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
)

type Page interface {
	// 是否有权限可见
	AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool
	VueRouter() Router
	HttpHandles(router gin.IRouter)
	Title() string     // 标题
	// 导航栏是否显示
	DisplayInNavigation(ctx *gin.Context, user interface{}) bool
}
