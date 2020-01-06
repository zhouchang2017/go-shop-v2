package filters

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/filters"
)

type InventoryStatusFilter struct {
	*filters.BooleanFilter
}



func NewInventoryStatusFilter() *InventoryStatusFilter {
	return &InventoryStatusFilter{BooleanFilter: filters.NewBooleanFilter()}
}

func (this InventoryStatusFilter) Apply(ctx *gin.Context, value interface{}, request *request.IndexRequest) error {
	spew.Dump(value)
	return nil
}

func (this InventoryStatusFilter) Key() string {
	return "status"
}

func (this InventoryStatusFilter) Name() string {
	return "状态"
}

func (this InventoryStatusFilter) DefaultValue(ctx *gin.Context) interface{} {
	return []interface{}{}
}

func (this InventoryStatusFilter) Options(ctx *gin.Context) []contracts.FilterOption {
	return []contracts.FilterOption{
		filters.NewSelectOption("等待处理", models.ITEM_PENDING),
		filters.NewSelectOption("锁定", models.ITEM_LOCKED),
		filters.NewSelectOption("良品", models.ITEM_OK),
		filters.NewSelectOption("不良品", models.ITEM_BAD),
	}
}
