package pages

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

var OrderLogisticsPage *orderLogisticsPage

// 物流详情页面
type orderLogisticsPage struct {
	srv      *services.OrderService
	trackRep *repositories.TrackRep
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
	router.GET("/tracks/:deliveryId/:wayBillId", func(ctx *gin.Context) {
		deliveryId := ctx.Param("deliveryId")
		if deliveryId == "" {
			err2.ErrorEncoder(nil, err2.Err422.F("查询物流失败 [deliveryId]is required"), ctx.Writer)
			return
		}
		wayBillId := ctx.Param("wayBillId")
		if wayBillId == "" {
			err2.ErrorEncoder(nil, err2.Err422.F("查询物流失败 [wayBillId]is required"), ctx.Writer)
			return
		}
		one, err := o.trackRep.FindOne(ctx, bson.M{
			"delivery_id": deliveryId,
			"way_bill_id": wayBillId,
		})
		if err != nil {
			err2.ErrorEncoder(nil, err, ctx.Writer)
			return
		}
		ctx.JSON(http.StatusOK, one)
	})
}

func (o orderLogisticsPage) Title() string {
	return "物流详情"
}

func (o orderLogisticsPage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}

func NewOrderLogisticsPage() *orderLogisticsPage {
	if OrderLogisticsPage == nil {
		OrderLogisticsPage = &orderLogisticsPage{srv: services.MakeOrderService(), trackRep: repositories.MakeTrackRep()}
	}
	return OrderLogisticsPage
}
