package charts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/vue/charts"
)

var PaymentLine *paymentLine

type paymentLine struct {
	*charts.Line
	srv *services.PaymentService
}

func (p paymentLine) Name() string {
	return "7日成交趋势"
}

func (p paymentLine) Columns() []string {
	return []string{"date", "total_amount", "count"}
}

func (p paymentLine) HttpHandle(ctx *gin.Context) (rows interface{}, err error) {
	end := utils.TodayEnd()
	return p.srv.RangePaymentCounts(ctx, end.AddDate(0, 0, -7), end)
}

func NewPaymentLine() *paymentLine {
	if PaymentLine == nil {
		PaymentLine = &paymentLine{
			Line: charts.NewLine(),
			srv:  services.MakePaymentService(),
		}
		PaymentLine.LabelMap(map[string]interface{}{
			"date":         "日期",
			"total_amount": "成交金额",
			"count":        "交易笔数",
		})
		PaymentLine.SetWidthFull()
		//PaymentLine.Dimension([]string{"date"})
		PaymentLine.XAxisTypeTime()
		PaymentLine.WithSettings("yAxisName", []string{"金额", "笔数"})
		PaymentLine.WithSettings("axisSite", map[string][]string{
			"right":[]string{"count"},
		})
	}
	return PaymentLine
}
