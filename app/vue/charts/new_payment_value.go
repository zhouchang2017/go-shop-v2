package charts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/vue/charts"
)

var NewPaymentValue *newPaymentValue

type newPaymentValue struct {
	*charts.Value
	srv *services.PaymentService
}

func NewNewPaymentValue() *newPaymentValue {
	if NewPaymentValue == nil {
		NewPaymentValue = &newPaymentValue{
			Value: charts.NewValue(),
			srv:   services.MakePaymentService(),
		}
	}
	return NewPaymentValue
}

func (v newPaymentValue) Columns() []string {
	return []string{}
}

func (v newPaymentValue) HttpHandle(ctx *gin.Context) (rows interface{}, err error) {
	count, err := v.srv.TodayPaymentCount(ctx)
	if err != nil {
		return 0, nil
	}
	return count, nil
}

func (v newPaymentValue) Name() string {
	return "当日付款金额"
}

func (this newPaymentValue) Component() string {
	return "cards-payment-value"
}
