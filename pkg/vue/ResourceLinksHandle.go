package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	"net/http"
)

func (this *ResourceWarp) resourceLinksIndexHandle(router gin.IRouter) {
	router.GET(fmt.Sprintf("links/%s", this.UriKey()), func(c *gin.Context) {
		res := []LensResponse{}
		links := this.resource.Links()
		for _, lens := range links {
			if lens.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
				res = append(res, LensResponse{
					Title:      lens.Title(),
					RouterName: lens.RouterName(),
				})
			}
		}
		c.JSON(http.StatusOK, res)
	})
}
