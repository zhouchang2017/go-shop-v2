package vue

import (
	"github.com/gin-gonic/gin"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"net/http"
	"reflect"
)

func (this *ResourceWarp) resourceIndexHandle(router gin.IRouter) {
	if r, ok := this.resource.(ResourceHttpIndex); ok && r.ResourceHttpIndex() {
		router.GET(this.UriKey(), func(c *gin.Context) {
			this.resource.Model()
			// 验证权限
			if !this.AuthorizedToViewAny(c) {
				c.AbortWithStatus(403)
				return
			}

			// 处理函数
			filter := &request.IndexRequest{}
			if err := c.ShouldBind(filter); err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}
			// hook 自定义查询
			r.IndexQuery(c, filter)

			results := <-this.resource.Repository().Pagination(c, filter)
			if results.Error != nil {
				err2.ErrorEncoder(nil, results.Error, c.Writer)
				return
			}
			indexResources := []Metable{}
			if reflect.TypeOf(results.Result).Kind() == reflect.Slice {
				valueOf := reflect.ValueOf(results.Result)
				len := valueOf.Len()
				for i := 0; i < len; i++ {
					model := valueOf.Index(i).Interface()
					indexResources = append(indexResources, NewResourceWarp(this.resource.Make(model), this.root).SerializeForIndex(c))
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"pagination": results.Pagination,
				"data":       indexResources,
			})
		})
	}

}
