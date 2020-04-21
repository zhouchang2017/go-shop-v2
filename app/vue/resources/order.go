package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/app/vue/charts"
	fields2 "go-shop-v2/app/vue/fields"
	"go-shop-v2/app/vue/filters"
	"go-shop-v2/app/vue/pages"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
	"net/http"
)

type Order struct {
	core.AbstractResource
	srv       *services.OrderService
	refundSrv *services.RefundService
	model     interface{}
}

func (order *Order) SearchPlaceholder() string {
	return "请输入订单号"
}

func (order *Order) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return order.srv.FindById(ctx, id)
}

func (order *Order) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	req.SetSearchField("order_no")
	return order.srv.Pagination(ctx, req)
}

func (order *Order) Title() string {
	return "订单"
}

func (order *Order) Icon() string {
	return "icons-clipboard"
}

func (order *Order) Group() string {
	return "Order"
}

func (order *Order) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(fields.WithMeta("min-width", 150)),
			fields.NewTextField("订单号", "OrderNo", fields.WithMeta("min-width", 150)),

			fields.NewCurrencyField("订单金额", "OrderAmount"),
			fields.NewCurrencyField("应付金额", "ActualAmount"),

			panels.NewPanel("收货信息",
				fields.NewTextField("收货人", "UserAddress.ContactName", fields.SetShowOnIndex(false)),
				fields.NewTextField("联系方式", "UserAddress.ContactPhone", fields.SetShowOnIndex(false)),
				fields.NewAreaCascader("省/市/区", "UserAddress", fields.SetShowOnIndex(false)),
				fields.NewTextField("详细地址", "UserAddress.Addr", fields.SetShowOnIndex(false)),
			),

			panels.NewPanel("用户",
				fields.NewAvatar("头像", "User.Avatar", fields.SetShowOnIndex(false)).RoundedFull(),
				fields.NewTextField("用户", "User.Nickname"),
				fields.NewStatusField("性别", "User.Gender", fields.SetShowOnIndex(false)).WithOptions([]*fields.StatusOption{
					fields.NewStatusOption("未知", 0),
					fields.NewStatusOption("男", 1),
					fields.NewStatusOption("女", 2),
				}),
			),

			panels.NewPanel("支付信息",
				fields.NewTextField("支付单号", "Payment.PaymentNo", fields.SetShowOnIndex(false)),
				fields.NewTextField("支付平台", "Payment.Platform", fields.SetShowOnIndex(false)),
				fields.NewCurrencyField("支付金额", "Payment.Amount", fields.SetShowOnIndex(false)),
				fields.NewDateTime("创建时间", "Payment.CreatedAt", fields.SetShowOnIndex(false)),
				fields.NewDateTime("支付时间", "Payment.PaymentAt", fields.SetShowOnIndex(false)),
			),
			//fields.NewHasManyField()
			fields2.NewOrderItemsField(),
			fields.NewDateTime("创建时间", "CreatedAt", fields.SetShowOnIndex(false)),
			fields.NewDateTime("更新时间", "UpdatedAt", fields.WithMeta("min-width", 150)),
			fields2.NewOrderStatusField(),

			fields.NewHasManyField("评论", &Comment{}),
		}
	}

}

func (order *Order) Model() interface{} {
	return order.model
}

func (order *Order) Make(mode interface{}) contracts.Resource {
	return &Order{
		srv:       order.srv,
		model:     mode,
		refundSrv: order.refundSrv,
	}
}

func (order *Order) SetModel(model interface{}) {
	order.model = model
}

func NewOrderResource() *Order {
	return &Order{srv: services.MakeOrderService(), model: &models.Order{}, refundSrv: services.MakeRefundService()}
}

func (this *Order) CustomHttpHandle(router gin.IRouter) {
	// 关闭订单
	router.PUT("/api/orders/:Order/cancel", func(ctx *gin.Context) {
		type cancelForm struct {
			Reason string `json:"reason"`
		}
		id := ctx.Param("Order")
		var form cancelForm
		reason := "暂时缺货"
		if err := ctx.ShouldBind(&form); err == nil {
			reason = form.Reason
		}
		if id == "" {
			err2.ErrorEncoder(ctx, err2.Err422.F("订单号异常"), ctx.Writer)
			return
		}

		order, err := this.srv.FindById(ctx, id)
		if err != nil {
			err2.ErrorEncoder(ctx, err2.Err422.F("订单不存在"), ctx.Writer)
			return
		}

		updatedOrder, err := this.srv.Cancel(ctx, order, reason)
		if err != nil {
			err2.ErrorEncoder(ctx, err, ctx.Writer)
			return
		}
		// 管理员取消订单事件推送
		rabbitmq.Dispatch(events.NewOrderClosedByAdminEvent(updatedOrder))
		ctx.JSON(http.StatusNoContent, nil)
	})

}

func (this *Order) Filters(ctx *gin.Context) []contracts.Filter {
	return []contracts.Filter{
		filters.NewOrderStatusFilter(),
	}
}

func (this Order) Cards(ctx *gin.Context) []contracts.Card {
	return []contracts.Card{
		charts.NewCountOrderPrePayValue(),
		charts.NewCountOrderPreSendValue(),
		charts.NewCountOrderRefundingValue(),
	}
}

func (this Order) Pages() []contracts.Page {
	return []contracts.Page{
		pages.NewOrderItemAggregatePage(),
	}
}
