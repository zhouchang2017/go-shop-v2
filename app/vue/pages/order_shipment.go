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
	srv          *services.OrderService
	logisticsSrv *services.LogisticsService
}

func NewOrderShipmentPage() *orderShipmentPage {
	return &orderShipmentPage{srv: services.MakeOrderService(), logisticsSrv: services.MakeLogisticsService()}
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
	// 获取微信小程序可用物流
	router.GET("mp/logistics", func(ctx *gin.Context) {
		delivery, _ := o.logisticsSrv.GetAllDelivery()
		ctx.JSON(http.StatusOK, delivery)
	})
	// 小程序物流助手下单
	router.POST("mp/logistics", func(ctx *gin.Context) {
		var form services.CreateExpressOrderOption
		if err := ctx.ShouldBind(&form); err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}
		if err := form.IsValid(); err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}
		order, err := o.logisticsSrv.AddOrder(ctx, form)
		if err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}
		ctx.JSON(http.StatusOK, order)
	})
	// 查询小程序物流信息
	router.GET("mp/logistics/get-order", func(ctx *gin.Context) {
		var form services.GetOrderOption
		if err := ctx.ShouldBind(&form); err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}
		order, err := o.logisticsSrv.GetOrder(ctx, &form)
		if err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}
		ctx.JSON(http.StatusOK, order)
	})
	// 取消物流助手下单
	router.POST("mp/logistics/cancel-order", func(ctx *gin.Context) {
		var form services.CancelOrderOption
		if err := ctx.ShouldBind(&form); err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}

		err := o.logisticsSrv.CancelExpressOrder(ctx, &form)
		if err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}
		ctx.JSON(http.StatusNoContent, nil)
	})
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
