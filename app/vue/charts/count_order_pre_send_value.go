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

var CountOrderPreSendValue *countOrderPreSendValue

type countOrderPreSendValue struct {
	*charts.Value
	srv *services.OrderService
}

func NewCountOrderPreSendValue() *countOrderPreSendValue {
	if CountOrderPreSendValue == nil {
		CountOrderPreSendValue = &countOrderPreSendValue{
			Value: charts.NewValue(),
			srv:   services.MakeOrderService(),
		}
	}
	return CountOrderPreSendValue
}

func (v countOrderPreSendValue) Columns() []string {
	return []string{}
}

func (v countOrderPreSendValue) HttpHandle(ctx *gin.Context) (rows interface{}, err error) {
	count, err := v.srv.CountByStatus(ctx, models.OrderStatusPreSend)
	if err != nil {
		return 0, nil
	}
	return count, nil
}

func (countOrderPreSendValue) Name() string {
	return "待发货订单"
}

func (v countOrderPreSendValue) Link() contracts.VueRouterOption {
	option := core.NewVueRouterOption("orders.index")
	filter := map[string]interface{}{
		"status": []int{models.OrderStatusPreSend},
	}
	if filterString, err := request.StructToFilterString(filter); err == nil {
		option.SetQuery(map[string]interface{}{
			"orders_filter": filterString,
		})
	}

	return option
}

func (v countOrderPreSendValue) Refresh() time.Duration {
	return time.Second * 10
}
