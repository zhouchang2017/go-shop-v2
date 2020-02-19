package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/fields"
)

type InventoryLog struct {
	rep   *repositories.InventoryLogRep
	model interface{}
}

func NewInventoryLogResource() *InventoryLog {
	return &InventoryLog{
		rep: repositories.NewInventoryLogRep(mongodb.GetConFn()),
	}
}

func (this *InventoryLog) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	inventoryId:= ctx.Query("inventory_id")
	if inventoryId!=""{
		req.AppendFilter("inventory_id",inventoryId)
	}
	results := <-this.rep.Pagination(ctx, req)
	return results.Result, results.Pagination, results.Error
}

func (InventoryLog) Title() string {
	return "库存操作日志"
}

func (InventoryLog) Group() string {
	return "Shop"
}

func (*InventoryLog) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}

func (*InventoryLog) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (*InventoryLog) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (*InventoryLog) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return false
}

func (*InventoryLog) Policy() interface{} {
	return nil
}

func (*InventoryLog) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("Before Qty", "BeforeQty"),
			fields.NewTextField("After Qty", "AfterQty"),
			fields.NewTextField("数量", "Name"),
			fields.NewTextField("操作用户", "User.Nickname"),
			fields.NewDateTime("操作时间", "UpdatedAt"),
		}
	}
}

func (this *InventoryLog) Model() interface{} {
	return this.model
}

func (this *InventoryLog) Make(mode interface{}) contracts.Resource {
	return &InventoryLog{
		model: mode,
		rep:   this.rep,
	}
}

func (this *InventoryLog) SetModel(model interface{}) {
	this.model = model
}

func (*InventoryLog) Lenses() []contracts.Lens {
	return []contracts.Lens{}
}

func (*InventoryLog) Pages() []contracts.Page {
	return []contracts.Page{}
}

func (*InventoryLog) Filters(ctx *gin.Context) []contracts.Filter {
	return []contracts.Filter{}
}

func (*InventoryLog) Actions(ctx *gin.Context) []contracts.Action {
	return []contracts.Action{}
}

func (*InventoryLog) Cards(ctx *gin.Context) []contracts.Card {
	return []contracts.Card{}
}
