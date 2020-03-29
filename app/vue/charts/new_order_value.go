package charts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/vue/charts"
)

var NewOrderValue *newOrderValue

type newOrderValue struct {
	*charts.Value
	srv *services.OrderService
}

func NewNewOrderValue() *newOrderValue {
	if NewOrderValue == nil {
		NewOrderValue = &newOrderValue{
			Value: charts.NewValue(),
			srv:   services.MakeOrderService(),
		}
	}
	return NewOrderValue
}

func (v newOrderValue) Columns() []string {
	return []string{}
}

func (v newOrderValue) HttpHandle(ctx *gin.Context) (rows interface{}, err error) {
	count := v.srv.TodayNewOrderCount(ctx)
	return count, nil
}

func (newOrderValue) Name() string {
	return "当日新增订单"
}
