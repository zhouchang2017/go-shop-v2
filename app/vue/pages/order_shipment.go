package pages

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"net/http"
)

var OrderShipmentPage *orderShipmentPage

// 订单发货页面
type orderShipmentPage struct {
	srv *services.OrderService
}

func NewOrderShipmentPage() *orderShipmentPage {
	return &orderShipmentPage{srv: services.MakeOrderService()}
}

func (o orderShipmentPage) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (o *orderShipmentPage) VueRouter() contracts.Router {
	router := core.NewRouter()
	router.RouterPath = "orders/:id/shipment"
	router.Name = "orders.shipment"
	router.RouterComponent = "orders/Shipment"
	router.Hidden = true
	router.WithMeta("ResourceName", "orders")
	router.WithMeta("Title", o.Title())
	return router
}

func (o *orderShipmentPage) HttpHandles(router gin.IRouter) {
	// 发货
	router.POST("orders/:Order/shipment", func(ctx *gin.Context) {
		id := ctx.Param("Order")
		if id == "" {
			err2.ErrorEncoder(nil, err2.Err422.F("缺少 order id 参数"), ctx.Writer)
			return
		}
		var form services.DeliverOption
		if err := ctx.ShouldBind(&form); err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}

		order, err := o.srv.FindById(ctx, id)
		if err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}

		model, err := o.srv.Deliver(ctx, order, &form)
		if err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}
		ctx.JSON(http.StatusOK, model)
	})
}

func (o orderShipmentPage) Title() string {
	return "发货"
}

func (o orderShipmentPage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}
