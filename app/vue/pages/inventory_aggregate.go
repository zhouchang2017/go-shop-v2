package pages

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"net/http"
)


var InventoryAggratePage *inventoryAggregate

// 自定义聚合页
type inventoryAggregate struct {
	service *services.InventoryService
	shopSrv *services.ShopService
	router  contracts.Router
}

func NewInventoryAggregatePage() *inventoryAggregate {
	if InventoryAggratePage == nil {
		InventoryAggratePage = &inventoryAggregate{
			service: services.MakeInventoryService(),
			shopSrv: services.MakeShopService(),
		}
	}
	return InventoryAggratePage
}

func (this *inventoryAggregate) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (this *inventoryAggregate) getShops() []*models.AssociatedShop {
	return this.shopSrv.GetAllAssociatedShops(context.Background())
}

func (this *inventoryAggregate) VueRouter() contracts.Router {
	if this.router == nil {
		router := core.NewRouter()
		router.RouterPath = "inventories/aggregate"
		router.Name = "inventories.aggregate"
		router.RouterComponent = "inventories/Aggregate"
		router.Hidden = true
		router.WithMeta("ResourceName", "inventories")
		router.WithMeta("Title", this.Title())
		router.WithMeta("shops", this.getShops())
		this.router = router
	}
	return this.router
}

func (this *inventoryAggregate) UriKey() string {
	return "inventories/aggregate"
}

func (this *inventoryAggregate) RouterName() string {
	return "inventories.aggregate"
}

func (this *inventoryAggregate) Component() string {
	return "inventories/Aggregate"
}

func (this *inventoryAggregate) HttpHandles(router gin.IRouter) {
	router.GET("aggregate/inventories", func(c *gin.Context) {
		// 验证权限
		if !this.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			c.AbortWithStatus(403)
			return
		}

		// 处理函数
		filter := &request.IndexRequest{}
		if err := c.ShouldBind(filter); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		data, pagination, err := this.service.Aggregate(c, filter)

		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pagination": pagination,
			"data":       data,
		})
	})
}

func (this *inventoryAggregate) Title() string {
	return "多门店聚合"
}

func (this *inventoryAggregate) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}
