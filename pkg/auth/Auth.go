package auth

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"net/http"
)

// 鉴权中间件
func AuthMiddleware(guard string) gin.HandlerFunc {
	return func(c *gin.Context) {
		statefulGuard, err := Auth.Guard(guard)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		ctx.GinCtxWithGuard(c, statefulGuard)
		statefulGuard.SetContext(c)

		user, err := statefulGuard.User()
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			//c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			c.Abort()
			return
		}
		ctx.GinCtxWithUser(c, user)
		c.Next()
	}
}
