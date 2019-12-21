package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"net/http"
)

func (this *ResourceWarp) resourceDeleteHandle(router gin.IRouter) {
	if r, ok := this.resource.(ResourceHttpDelete); ok && r.ResourceHttpDelete() {
		router.DELETE(fmt.Sprintf("%s/:%s", this.UriKey(), this.SingularLabel()), func(c *gin.Context) {

			// 验证权限
			if ok, _ := this.AuthorizedToDelete(c); !ok {
				c.AbortWithStatus(403)
				return
			}

			// 处理函数
			err := <-this.resource.Repository().Delete(c, c.Param(this.SingularLabel()))
			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}
			go this.resource.Deleted(c,c.Param(this.SingularLabel()))
			c.JSON(http.StatusOK, nil)
		})
	}

}

func (this *ResourceWarp) resourceForceDeleteHandle(router gin.IRouter) {
	if r, ok := this.resource.(ResourceHttpForceDelete); ok && r.ResourceHttpForceDelete() {
		router.DELETE(fmt.Sprintf("%s/:%s/force", this.UriKey(), this.SingularLabel()), func(c *gin.Context) {

			// 验证权限
			if ok, _ := this.AuthorizedToForceDelete(c); !ok {
				c.AbortWithStatus(403)
				return
			}

			// 处理函数
			context := ctx.WithForce(c, true)
			err := <-this.resource.Repository().Delete(context, c.Param(this.SingularLabel()))
			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}
			// force deleted hook
			go this.resource.ForceDeleted(c,c.Param(this.SingularLabel()))
			c.JSON(http.StatusOK, nil)
		})
	}

}
