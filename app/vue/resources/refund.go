package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
)

type Refund struct {
	core.AbstractResource
	srv   *services.RefundService
	model interface{}
}

func (r Refund) SearchPlaceholder() string {
	return "请输入退款单号"
}

func NewRefundResource() *Refund {
	return &Refund{
		srv:   services.MakeRefundService(),
		model: &models.Refund{},
	}
}

func (r Refund) Title() string {
	return "退款"
}

func (r *Refund) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	req.SetSearchField("refund_no")
	return r.srv.Pagination(ctx, req)
}

func (r Refund) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {

		return []interface{}{
			fields.NewIDField(fields.WithMeta("min-width", 150)),
			fields.NewTextField("退款单号", "RefundNo", fields.WithMeta("min-width", 150)),
			fields.NewTextField("订单号", "OrderNo", fields.WithMeta("min-width", 150)).Link(&Order{},"OrderId"),
			fields.NewTextField("支付单号", "PaymentNo", fields.WithMeta("min-width", 150)),

			fields.NewCurrencyField("退款金额", "TotalAmount"),

			fields.NewStatusField("回调状态", "ReturnCode").WithOptions([]*fields.StatusOption{
				fields.NewStatusOption("成功", "SUCCESS").Success(),
				fields.NewStatusOption("失败", "FAIL").Error(),
				fields.NewStatusOption("N/A", "").Cancel(),
			}),

			fields.NewStatusField("状态", "Status").WithOptions([]*fields.StatusOption{
				fields.NewStatusOption("申请退款", models.RefundStatusApply).Error(),
				fields.NewStatusOption("同意退款", models.RefundStatusAgreed).Info(),
				fields.NewStatusOption("拒绝退款", models.RefundStatusReject).Cancel(),
				fields.NewStatusOption("退款中", models.RefundStatusRefunding).Warning(),
				fields.NewStatusOption("退款完成", models.RefundStatusDone).Success(),
				fields.NewStatusOption("退款关闭", models.RefundStatusClosed).Cancel(),
			}),


			fields.NewDateTime("创建时间", "CreatedAt", fields.SetShowOnIndex(false)),
			fields.NewDateTime("更新时间", "UpdatedAt", fields.WithMeta("min-width", 150)),
		}
	}
}

func (r *Refund) Model() interface{} {
	return r.model
}

func (r *Refund) Make(mode interface{}) contracts.Resource {
	return &Refund{
		AbstractResource: core.AbstractResource{},
		srv:              r.srv,
		model:            mode,
	}
}

func (r *Refund) SetModel(model interface{}) {
	r.model = model
}

func (r Refund) Icon() string {
	return "icons-currency-dollar"
}

func (r Refund) Group() string {
	return "Order"
}
