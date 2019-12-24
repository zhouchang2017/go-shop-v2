package vue

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
)

type Link interface {
	Title() string
	RouterName() string
	AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool // 授权
}
