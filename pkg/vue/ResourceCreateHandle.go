package vue

import (
	"github.com/gin-gonic/gin"
	err2 "go-shop-v2/pkg/err"
	"net/http"
)

func (this *ResourceWarp) resourceCreateHandle(router gin.IRouter) {
	if creatable, ok := this.resource.(ResourceHttpCreate); ok && creatable.ResourceHttpCreate() {
		router.POST(this.UriKey(), func(c *gin.Context) {
			// 验证权限
			if !this.AuthorizedToCreate(c) {
				c.AbortWithStatus(403)
				return
			}

			// 资源处理表单
			entity, err := creatable.CreateFormParse(c)
			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}

			results := <-this.resource.Repository().Create(c, entity)
			if results.Error != nil {
				err2.ErrorEncoder(nil, results.Error, c.Writer)
				return
			}
			// created hook
			go this.resource.Created(c,results.Result)

			c.JSON(http.StatusCreated, gin.H{"id": results.Id})
		})
	}

}
