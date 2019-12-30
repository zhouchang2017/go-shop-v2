package contracts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
)

// 列表api权限接口
type AuthorizedToViewAny interface {
	ViewAny(ctx *gin.Context, user auth.Authenticatable) bool
}

// 详情api权限接口
type AuthorizedToView interface {
	View(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

// 创建api权限接口
type AuthorizedToCreate interface {
	Create(ctx *gin.Context, user auth.Authenticatable) bool
}

// 更新api权限接口
type AuthorizedToUpdate interface {
	Update(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

// 删除api权限接口
type AuthorizedToDelete interface {
	Delete(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

// 硬删除api权限接口
type AuthorizedToForceDelete interface {
	ForceDelete(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

// 还原api权限接口
type AuthorizedToRestore interface {
	Restore(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}