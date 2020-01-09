package charts

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/vue/charts"
)

var ShopsInventoryBar *shopsInventoryBar

type shopsInventoryBar struct {
	*charts.Bar
	rep *repositories.InventoryRep
}

func NewShopsInventoryBar() *shopsInventoryBar {
	if ShopsInventoryBar == nil {
		ShopsInventoryBar = &shopsInventoryBar{
			Bar: charts.NewBar(),
			rep: repositories.NewInventoryRep(mongodb.GetConFn()),
		}
		ShopsInventoryBar.LabelMap(map[string]interface{}{
			"shop_name": "门店名称",
			"total":     "总计",
			"status_0":  "待确认",
			"status_1":  "锁定",
			"status_2":  "良品",
			"status_3":  "不良品",
		})
		ShopsInventoryBar.SetWidth50Percent()
		ShopsInventoryBar.Stack([]string{"status_0", "status_1", "status_2", "status_3"})
	}

	return ShopsInventoryBar
}

func (shopsInventoryBar) Name() string {
	return "库存统计"
}

func (shopsInventoryBar) Columns() []string {
	return []string{"shop_name", "total", "status_2", "status_0", "status_1", "status_3"}
}

func (this shopsInventoryBar) HttpHandle(ctx *gin.Context) (rows interface{}, err error) {
	//shopId := ctx.Param("resourceId")
	data, err := this.rep.AggregateStockByShops(ctx)
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
		}
		res = append(res, statusLine)
	}
	// "shop_name":"shopA","total":100,
	return res, nil
}
