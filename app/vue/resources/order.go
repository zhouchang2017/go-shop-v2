package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
)

func init() {
	register(NewOrderResource)
}

type Order struct {
	core.AbstractResource
	rep *repositories.OrderRep
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
	return "icons-store"
}

func (order *Order) Group() string {
	return "Order"
}

func (order *Order) Repository() repository.IRepository {
	return order.rep
}

func (order *Order) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("订单号", "OrderNo"),
			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),
			// todo
		}
	}

}

func (order *Order) Model() interface{} {
	return nil
}

func (order *Order) Make(mode interface{}) contracts.Resource {
	return &Order{
		rep: order.rep,
	}
}

func (order *Order) SetModel(model interface{}) {
	panic("implement me")
}

func NewOrderResource(rep *repositories.OrderRep) *Order {
	return &Order{rep: rep}
}
