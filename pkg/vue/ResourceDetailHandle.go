package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"net/http"
)

func (this *ResourceWarp) resourceDetailHandle(router gin.IRouter) {
	if r, ok := this.resource.(ResourceHttpShow); ok && r.ResourceHttpShow() {
		router.GET(fmt.Sprintf("%s/:%s", this.UriKey(), this.SingularLabel()), func(c *gin.Context) {

			ctx.GinCtxWithTrashed(c)
			// 验证权限
			ok, model := this.AuthorizedToView(c);
			if !ok {
				c.AbortWithStatus(403)
				return
			}
			if extend, ok := this.resource.(CustomResourceHttpShow); ok {
				resModel, err := extend.CustomResourceHttpShow(c, c.Param(this.SingularLabel()))
				if err != nil {
					err2.ErrorEncoder(nil, err, c.Writer)
					return
				}
				model = resModel
			} else {
				if model == nil {
					result := <-this.resource.Repository().FindById(c, c.Param(this.SingularLabel()))
					if result.Error != nil {
						err2.ErrorEncoder(nil, result.Error, c.Writer)
						return
					}
					model = result.Result
				}
			}
			resource:= this.resource.Make(model)
			resource.SetRoot(this.root)

			c.JSON(http.StatusOK, NewResourceWarp(resource, this.root).SerializeForDetail(c))
		})
	}

}
