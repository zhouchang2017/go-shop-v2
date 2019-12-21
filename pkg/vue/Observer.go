package vue

import "github.com/gin-gonic/gin"

type Observer interface {
	Created(ctx *gin.Context, resource interface{})
	Updated(ctx *gin.Context, resource interface{})
	Deleted(ctx *gin.Context, id string)
	Restored(ctx *gin.Context,resource interface{})
	ForceDeleted(ctx *gin.Context,id string)
}
