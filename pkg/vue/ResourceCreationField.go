package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (this *ResourceWarp) resourceCreationFieldHandle(router gin.IRouter) {
	if creatable, ok := this.resource.(ResourceHttpCreate); ok && creatable.ResourceHttpCreate() {
		router.GET(fmt.Sprintf("creation-fields/%s", this.UriKey()), func(c *gin.Context) {
			// 验证权限
			if !this.AuthorizedToCreate(c) {
				c.AbortWithStatus(403)
				return
			}

			fields, panels := this.resolveCreationFields(c)
			c.JSON(http.StatusOK, gin.H{
				"fields": fields,
				"panels": panels,
			})
		})
	}
}
