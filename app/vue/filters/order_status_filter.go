package filters

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/filters"
	"go.mongodb.org/mongo-driver/bson"
)

type OrderStatusFilter struct {
	*filters.BooleanFilter
}

func NewOrderStatusFilter() *OrderStatusFilter {
	return &OrderStatusFilter{BooleanFilter: filters.NewBooleanFilter()}
}

func (this OrderStatusFilter) Apply(ctx *gin.Context, value interface{}, request *request.IndexRequest) error {
	status := value.([]interface{})
	if len(status) > 0 {
		request.AppendFilter("status", bson.M{"$in": value})
	}
	return nil
}

func (this OrderStatusFilter) Key() string {
	return "status"
}

func (this OrderStatusFilter) Name() string {
	return "状态"
}

func (this OrderStatusFilter) DefaultValue(ctx *gin.Context) interface{} {
	return []interface{}{}
}

func (this OrderStatusFilter) Options(ctx *gin.Context) []contracts.FilterOption {
	return []contracts.FilterOption{
		filters.NewSelectOption("已取消", models.OrderStatusFailed),
		filters.NewSelectOption("待付款", models.OrderStatusPrePay),
		filters.NewSelectOption("支付成功", models.OrderStatusPaid),
		filters.NewSelectOption("待发货", models.OrderStatusPreSend),
		filters.NewSelectOption("待收货", models.OrderStatusPreConfirm),
		filters.NewSelectOption("已完成", models.OrderStatusDone),
		filters.NewSelectOption("待评价", models.OrderStatusPreEvaluate),
	}
}
