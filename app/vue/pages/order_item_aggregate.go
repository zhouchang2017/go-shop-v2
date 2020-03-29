package pages

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"net/http"
)

var OrderItemAggregate *orderItemAggregate

type orderItemAggregate struct {
	srv    *services.OrderService
	router contracts.Router
}

func (o orderItemAggregate) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (o *orderItemAggregate) VueRouter() contracts.Router {
	if o.router == nil {
		router := core.NewRouter()
		router.RouterPath = "orders/aggregate/unit"
		router.Name = "orders.aggregate.unit"
		router.RouterComponent = "orders/UnitAggregate"
		router.Hidden = true
		router.WithMeta("ResourceName", "orders")
		router.WithMeta("Title", o.Title())
		o.router = router
	}
	return o.router
}

func (o orderItemAggregate) HttpHandles(router gin.IRouter) {
	router.GET("aggregate/orders/unit", func(c *gin.Context) {
		// 验证权限
		if !o.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			c.AbortWithStatus(403)
			return
		}

		// 处理函数
		filter := &request.IndexRequest{}
		if err := c.ShouldBind(filter); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		data, pagination, err := o.srv.AggregateOrderItem(c, filter)

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

func (o orderItemAggregate) Title() string {
	return "门店销售聚合"
}

func (o orderItemAggregate) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}

func NewOrderItemAggregatePage() *orderItemAggregate {
	if OrderItemAggregate == nil {
		OrderItemAggregate = &orderItemAggregate{
			srv: services.MakeOrderService(),
		}
	}
	return OrderItemAggregate
}
