package charts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/charts"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"time"
)

// 退款中的订单统计
var CountOrderRefundingValue *countOrderRefundingValue

type countOrderRefundingValue struct {
	*charts.Value
	srv *services.OrderService
}

func NewCountOrderRefundingValue() *countOrderRefundingValue {
	if CountOrderRefundingValue == nil {
		CountOrderRefundingValue = &countOrderRefundingValue{
			Value: charts.NewValue(),
			srv:   services.MakeOrderService(),
		}
	}
	return CountOrderRefundingValue
}

func (v countOrderRefundingValue) Columns() []string {
	return []string{}
}

func (v countOrderRefundingValue) HttpHandle(ctx *gin.Context) (rows interface{}, err error) {
	count := v.srv.RefundingCount(ctx)
	if err != nil {
		return 0, nil
	}
	return count, nil
}

func (countOrderRefundingValue) Name() string {
	return "退款订单"
}

func (v countOrderRefundingValue) Link() contracts.VueRouterOption {
	option := core.NewVueRouterOption("orders.index")
	filter := map[string]interface{}{
		"status": []string{"RefundStatusApply", "RefundStatusAgreed", "RefundStatusRefunding"},
	}
	if filterString, err := request.StructToFilterString(filter); err == nil {
		option.SetQuery(map[string]interface{}{
			"orders_filter": filterString,
		})
	}

	return option
}

func (v countOrderRefundingValue) Refresh() time.Duration {
	return time.Second * 10
}
