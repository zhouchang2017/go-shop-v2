package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/err"
	"net/http"
)

func (this *ResourceWarp) resourceRestoreHandle(router gin.IRouter) {
	if r, ok := this.resource.(ResourceHttpRestore); ok && r.ResourceHttpRestore() {
		router.PUT(fmt.Sprintf("%s/:%s/restore", this.UriKey(), this.SingularLabel()), func(c *gin.Context) {

			// 验证权限
			if ok, _ := this.AuthorizedToRestore(c); !ok {
				c.AbortWithStatus(403)
				return
			}

			// 处理函数
			result := <-this.resource.Repository().Restore(c, c.Param(this.SingularLabel()))
			if result.Error != nil {
				err.ErrorEncoder(nil, result.Error, c.Writer)
				return
			}
			// restored hook
			go this.resource.Restored(c, result.Result)
			c.JSON(http.StatusOK, nil)
		})
	}

}
