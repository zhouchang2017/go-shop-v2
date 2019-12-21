package vue

import "github.com/gin-gonic/gin"

// 列表页路由生成监听
type ListenerIndexRouteCreated interface {
	OnIndexRouteCreated(ctx *gin.Context, router *Router)
}

// 详情页路由生成监听
type ListenerDetailRouteCreated interface {
	OnDetailRouteCreated(ctx *gin.Context, router *Router)
}

// 更新页路由生成监听
type ListenerUpdateRouteCreated interface {
	OnUpdateRouteCreated(ctx *gin.Context, router *Router)
}

// 创建页路由生成监听
type ListenerCreateRouteCreated interface {
	OnCreateRouteCreated(ctx *gin.Context, router *Router)
}
