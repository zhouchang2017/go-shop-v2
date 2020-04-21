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
		var orderStatus []float64
		var refundStatus []string
		for _, s := range status {
			if order, ok := s.(float64); ok {
				orderStatus = append(orderStatus, order)
				continue
			}
			if refund, ok := s.(string); ok {
				refundStatus = append(refundStatus, refund)
			}
		}
		if len(orderStatus) > 0 {
			request.AppendFilter("status", bson.M{"$in": orderStatus})
		}
		if len(refundStatus) > 0 {
			var s []int
			for _, i := range refundStatus {
				s = append(s, refundStatusToInt(i))
			}
			request.AppendFilter("refunds", bson.M{"$elemMatch": bson.M{"status": bson.M{"$in": s}}})
		}
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
		filters.NewSelectOption("已收货", models.OrderStatusConfirm),
		filters.NewSelectOption("待评价", models.OrderStatusPreEvaluate),
		filters.NewSelectOption("交易完成", models.OrderStatusDone),
		filters.NewSelectOption("退款申请", "RefundStatusApply"),
		filters.NewSelectOption("同意退款", "RefundStatusAgreed"),
		filters.NewSelectOption("拒绝退款", "RefundStatusReject"),
		filters.NewSelectOption("退款中", "RefundStatusRefunding"),
		filters.NewSelectOption("退款完成", "RefundStatusDone"),
		filters.NewSelectOption("退款关闭", "RefundStatusClosed"),
		//filters.NewSelectOption("订单申请退款", models.OrderStatusRefundApply),
		//filters.NewSelectOption("同意退款", models.OrderStatusRefundAgreed),
		//filters.NewSelectOption("拒绝退款", models.OrderStatusRefundReject),
		//filters.NewSelectOption("退款中", models.OrderStatusRefunding),
		//filters.NewSelectOption("退款已完成", models.OrderStatusRefundDone),
	}
}

func refundStatusToInt(s string) int {
	switch s {
	case "RefundStatusApply":
		return models.RefundStatusApply
	case "RefundStatusAgreed":
		return models.RefundStatusAgreed
	case "RefundStatusReject":
		return models.RefundStatusReject
	case "RefundStatusRefunding":
		return models.RefundStatusRefunding
	case "RefundStatusDone":
		return models.RefundStatusDone
	default:
		return models.RefundStatusClosed
	}
}
