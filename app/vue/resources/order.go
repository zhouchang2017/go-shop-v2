package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

type Order struct {
	core.AbstractResource
	rep *repositories.OrderRep
	model interface{}
}

func (order *Order) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	result := <-order.rep.FindById(ctx, id)
	return result.Result, result.Error
}

func (order *Order) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	results := <-order.rep.Pagination(ctx, req)
	return results.Result, results.Pagination, results.Error
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
			fields.NewIDField(),
			fields.NewTextField("订单号", "OrderNo"),
			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),
			fields.NewTextField("订单金额", "OrderAmount"),
			fields.NewTextField("实付金额", "ActualAmount"),
			fields.NewStatusField("订单状态", "Status").WithOptions([]*fields.StatusOption{
				fields.NewStatusOption("待支付", 0).Cancel(),
				fields.NewStatusOption("待发货", 1).Warning(),
				fields.NewStatusOption("待收货", 2).Info(),
				fields.NewStatusOption("已完成", 3).Success(),
				fields.NewStatusOption("已取消", 4).Danger(),
				fields.NewStatusOption("处理失败", 5).Error(),
			}),

			panels.NewPanel("收货信息",
				fields.NewTextField("收货人", "UserAddress.ContactName", fields.SetShowOnIndex(false)),
				fields.NewTextField("联系方式", "UserAddress.ContactPhone", fields.SetShowOnIndex(false)),
				fields.NewAreaCascader("省/市/区", "UserAddress"),
				fields.NewTextField("详细地址", "UserAddress.Addr", fields.SetShowOnIndex(false)),
			),

			panels.NewPanel("用户",
				fields.NewTextField("用户", "Nickname"),
				fields.NewTextField("头像", "Avatar", fields.SetShowOnIndex(false)),
				fields.NewStatusField("性别", "Gender", fields.SetShowOnIndex(false)).WithOptions([]*fields.StatusOption{
					fields.NewStatusOption("未知", 0),
					fields.NewStatusOption("男", 1),
					fields.NewStatusOption("女", 2),
				}),
			),

			fields.NewTable("物流信息", "Logistics", func() []contracts.Field {
				return []contracts.Field{
					fields.NewTextField("类型", "Enterprise"),
					fields.NewTextField("单号", "TrackNo"),
				}
			}),

			fields.NewTable("支付信息", "Payment", func() []contracts.Field {
				return []contracts.Field{
					fields.NewTextField("支付平台", "Platform"),
					fields.NewTextField("支付单号", "PaymentNo"),
				}
			}),
		}
	}

}

func (order *Order) Model() interface{} {
	return order.model
}

func (order *Order) Make(mode interface{}) contracts.Resource {
	return &Order{
		rep: order.rep,
	}
}

func (order *Order) SetModel(model interface{}) {
	order.model = model
}

func NewOrderResource() *Order {
	return &Order{rep: repositories.NewOrderRep(mongodb.GetConFn())}
}
