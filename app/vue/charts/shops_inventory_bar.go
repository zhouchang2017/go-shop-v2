package charts

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/vue/charts"
)

var ShopsInventoryBar *shopsInventoryBar

type shopsInventoryBar struct {
	*charts.Bar
	srv *services.InventoryService
}

func NewShopsInventoryBar() *shopsInventoryBar {
	if ShopsInventoryBar == nil {
		ShopsInventoryBar = &shopsInventoryBar{
			Bar: charts.NewBar(),
			srv: services.MakeInventoryService(),
		}
		ShopsInventoryBar.LabelMap(map[string]interface{}{
			"shop_name":  "门店名称",
			"total":      "总计",
			"qty":        "非锁定库存",
			"locked_qty": "锁定库存",
			"status_0":   "良品",
			"status_1":   "不良品",
		})
		ShopsInventoryBar.SetWidth50Percent()
		ShopsInventoryBar.Stack([]string{"status_0", "status_1"})
	}

	return ShopsInventoryBar
}

func (shopsInventoryBar) Name() string {
	return "库存统计"
}

func (shopsInventoryBar) Columns() []string {
	return []string{"shop_name", "total", "locked_qty", "status_0", "status_1", "qty"}
}

func (this shopsInventoryBar) HttpHandle(ctx *gin.Context) (rows interface{}, err error) {
	//shopId := ctx.Param("resourceId")
	data, err := this.srv.AggregateStockByShops(ctx)
	if err != nil {
		return
	}
	res := []interface{}{}
	for _, item := range data {
		statusLine := map[string]interface{}{}
		statusLine["shop_id"] = item.ShopId
		statusLine["shop_name"] = item.ShopName
		statusLine["total"] = item.Total
		for _, status := range item.Status {
			statusLine[fmt.Sprintf("status_%d", status.Status)] = status.Qty
			statusLine["qty"] = status.Qty
			statusLine["locked_qty"] = status.LockedQty
		}
		res = append(res, statusLine)
	}
	// "shop_name":"shopA","total":100,
	return res, nil
}
