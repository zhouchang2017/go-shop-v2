package pages

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
)

var OrderRefundPage *orderRefundPage

type orderRefundPage struct {
	srv *services.OrderService
}

func NewOrderRefundPage() *orderRefundPage {
	return &orderRefundPage{srv: services.MakeOrderService()}
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

}

func (o orderRefundPage) Title() string {
	return "退款详情"
}

func (o orderRefundPage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}
