package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"net/http"
)

func (this *ResourceWarp) resourceUpdateHandle(router gin.IRouter) {
	if upgradeable, ok := this.resource.(ResourceHttpUpdate); ok && upgradeable.ResourceHttpUpdate() {
		router.PUT(fmt.Sprintf("%s/:%s", this.UriKey(), this.SingularLabel()), func(c *gin.Context) {
			var result repository.QueryResult
			// 验证权限
			ok, model := this.AuthorizedToUpdate(c)
			if !ok {
				c.AbortWithStatus(403)
				return
			}
			if model == nil {
				result := <-this.resource.Repository().FindById(c, c.Param(this.SingularLabel()))
				if result.Error != nil {
					err2.ErrorEncoder(nil, result.Error, c.Writer)
					return
				}
				model = result.Result
			}

			// 表单处理
			entity, err := upgradeable.UpdateFormParse(c, model)
			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}
			// 处理函数
			queryResult := <-this.resource.Repository().Save(c, entity)
			if queryResult.Error != nil {
				err2.ErrorEncoder(nil, result.Error, c.Writer)
				return
			}
			// updated hook
			go this.resource.Updated(c, queryResult.Result)
			c.JSON(http.StatusOK, nil)
		})
	}

}
