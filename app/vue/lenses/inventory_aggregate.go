package lenses

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue"
	"net/http"
)

// 库存聚合
type InventoryAggregateLens struct {
	resource vue.Resource
	service  *services.InventoryService
	helper   *vue.ResourceHelper
}

func (this InventoryAggregateLens) UriKey() string {
	return "aggregate"
}

func NewInventoryAggregateLens(resource vue.Resource, service *services.InventoryService) vue.Lens {
	return &InventoryAggregateLens{
		resource: resource,
		service:  service,
		helper:   vue.NewResourceHelper(resource),
	}
}

func (InventoryAggregateLens) Title() string {
	return "多门店聚合"
}

func (this *InventoryAggregateLens) Component() string {
	return fmt.Sprintf(`%s/Aggregate`, this.helper.UriKey())
}

func (*InventoryAggregateLens) VueRouter() vue.Router {
	panic("implement me")
}

func (i *InventoryAggregateLens) HttpHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		filter := &request.IndexRequest{}
		if err := ctx.ShouldBind(filter); err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}

		data, pagination, err := i.service.Aggregate(ctx, filter)
		if err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"pagination": pagination,
			"data":       data,
		})
	}
}

func (*InventoryAggregateLens) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}
