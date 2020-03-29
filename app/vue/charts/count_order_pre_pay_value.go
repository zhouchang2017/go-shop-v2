package charts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/charts"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"time"
)

var CountOrderPrePayValue *countOrderPrePayValue

type countOrderPrePayValue struct {
	*charts.Value
	srv *services.OrderService
}

func (v countOrderPrePayValue) Refresh() time.Duration {
	return time.Second * 10
}

func NewCountOrderPrePayValue() *countOrderPrePayValue {
	if CountOrderPrePayValue == nil {
		CountOrderPrePayValue = &countOrderPrePayValue{
			Value: charts.NewValue(),
			srv:   services.MakeOrderService(),
		}
	}
	return CountOrderPrePayValue
}

func (v countOrderPrePayValue) Columns() []string {
	return []string{}
}

func (v countOrderPrePayValue) HttpHandle(ctx *gin.Context) (rows interface{}, err error) {
	count, err := v.srv.CountByStatus(ctx, models.OrderStatusPrePay)
	if err != nil {
		return 0, nil
	}
	return count, nil
}

func (countOrderPrePayValue) Name() string {
	return "待付款订单"
}

func (v countOrderPrePayValue) Link() contracts.VueRouterOption {
	option := core.NewVueRouterOption("orders.index")
	filter := map[string]interface{}{
		"status": []int{models.OrderStatusPrePay},
	}
	if filterString, err := request.StructToFilterString(filter); err == nil {
		option.SetQuery(map[string]interface{}{
			"orders_filter": filterString,
		})
	}

	return option
}
