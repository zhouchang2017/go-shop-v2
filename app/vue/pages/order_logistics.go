package pages

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
)

var OrderLogisticsPage *orderLogisticsPage

// 物流详情页面
type orderLogisticsPage struct {
	srv *services.OrderService
}

func (o orderLogisticsPage) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (o *orderLogisticsPage) VueRouter() contracts.Router {
	router := core.NewRouter()
	router.RouterPath = "orders/:id/logistics"
	router.Name = "orders.logistics"
	router.RouterComponent = "orders/Logistics"
	router.Hidden = true
	router.WithMeta("ResourceName", "orders")
	router.WithMeta("Title", o.Title())
	return router
}

func (o orderLogisticsPage) HttpHandles(router gin.IRouter) {

}

func (o orderLogisticsPage) Title() string {
	return "物流详情"
}

func (o orderLogisticsPage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}

func NewOrderLogisticsPage() *orderLogisticsPage {
	if OrderLogisticsPage == nil {
		OrderLogisticsPage = &orderLogisticsPage{srv: services.MakeOrderService()}
	}
	return OrderLogisticsPage
}
