package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	"net/http"
)

type LensResponse struct {
	Title      string `json:"title"`
	RouterName string `json:"router_name"`
}

// 资源对应自定义聚合路由
func (this *ResourceWarp) resourceLensesIndexHandle(router gin.IRouter) {
	router.GET(fmt.Sprintf("lenses/%s", this.UriKey()), func(c *gin.Context) {
		res := []LensResponse{}
		lenses := this.resource.Lenses()
		for _, lens := range lenses {
			if lens.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
				res = append(res, LensResponse{
					Title:      lens.Title(),
					RouterName: lensRouterName(lens, this.UriKey()),
				})
			}
		}
		c.JSON(http.StatusOK, res)
	})
}

// 资源对应自定义聚合路由处理
func (this *ResourceWarp) resourceLensesDetailHandle(router gin.IRouter) {
	for _, lens := range this.resource.Lenses() {

		router.GET(lensApiUri(lens, this.UriKey()), func(c *gin.Context) {
			// 权限中间件
			if !lens.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
				c.AbortWithStatus(403)
				return
			}
		}, lens.HttpHandle())

	}
}
