package pages

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"net/http"
)

var OrderRefundPage *orderRefundPage

type orderRefundPage struct {
	srv       *services.OrderService
	refundSrv *services.RefundService
}

func NewOrderRefundPage() *orderRefundPage {
	return &orderRefundPage{srv: services.MakeOrderService(), refundSrv: services.MakeRefundService()}
}

func (o orderRefundPage) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (o orderRefundPage) VueRouter() contracts.Router {
	router := core.NewRouter()
	router.RouterPath = "orders/:id/refund"
	router.Name = "orders.refund"
	router.RouterComponent = "orders/Refund"
	router.Hidden = true
	router.WithMeta("ResourceName", "orders")
	router.WithMeta("Title", o.Title())
	return router
}

func (o orderRefundPage) HttpHandles(router gin.IRouter) {
	// 获取订单关联退款
	router.GET("/api/orders/:Order/refunds", func(ctx *gin.Context) {
		orderId := ctx.Param("Order")
		if orderId == "" {
			err2.ErrorEncoder(ctx, err2.Err422.F("缺少order_id"), ctx.Writer)
			return
		}
		refunds, err := o.refundSrv.FindRefundByOrderId(ctx, orderId)
		if err != nil {
			err2.ErrorEncoder(ctx, err, ctx.Writer)
			return
		}
		ctx.JSON(http.StatusOK, refunds)
	})
	// 同意退款
	router.PUT("/api/refunds/:Refund/agree", func(ctx *gin.Context) {
		refundId := ctx.Param("Refund")
		if refundId == "" {
			err2.ErrorEncoder(ctx, err2.Err422.F("缺少refund_id"), ctx.Writer)
			return
		}

		refund, err := o.refundSrv.AgreeRefund(ctx, refundId)
		if err != nil {
			err2.ErrorEncoder(ctx, err, ctx.Writer)
			return
		}
		rabbitmq.Dispatch(events.NewOrderRefundChangeEvent(refund))
		ctx.JSON(http.StatusOK, refund)
	})
	// 拒绝退款
	router.PUT("/api/refunds/:Refund/reject", func(ctx *gin.Context) {
		refundId := ctx.Param("Refund")
		if refundId == "" {
			err2.ErrorEncoder(ctx, err2.Err422.F("缺少refund_id"), ctx.Writer)
			return
		}
		form := services.RefundOption{
			RefundId: refundId,
		}
		if err := ctx.ShouldBind(&form); err != nil {
			err2.ErrorEncoder(ctx, err, ctx.Writer)
			return
		}

		admin := ctx2.GetUser(ctx).(*models.Admin)

		refund, err := o.refundSrv.RejectRefund(ctx, admin, &form)
		if err != nil {
			err2.ErrorEncoder(ctx, err, ctx.Writer)
			return
		}
		rabbitmq.Dispatch(events.NewOrderRefundChangeEvent(refund))
		ctx.JSON(http.StatusOK, refund)
	})
	// 关闭退款
	router.PUT("/api/refunds/:Refund/cancel", func(ctx *gin.Context) {
		refundId := ctx.Param("Refund")
		if refundId == "" {
			err2.ErrorEncoder(ctx, err2.Err422.F("缺少refund_id"), ctx.Writer)
			return
		}
		form := services.RefundOption{
			RefundId: refundId,
		}
		if err := ctx.ShouldBind(&form); err != nil {
			err2.ErrorEncoder(ctx, err, ctx.Writer)
			return
		}
		admin := ctx2.GetUser(ctx).(*models.Admin)
		refund, err := o.refundSrv.CancelRefund(ctx, &form, admin, false)
		if err != nil {
			err2.ErrorEncoder(ctx, err, ctx.Writer)
			return
		}
		rabbitmq.Dispatch(events.NewOrderRefundChangeEvent(refund))
		ctx.JSON(http.StatusOK, refund)
	})
}

func (o orderRefundPage) Title() string {
	return "退款详情"
}

func (o orderRefundPage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}
