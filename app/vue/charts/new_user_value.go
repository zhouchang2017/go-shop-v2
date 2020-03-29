package charts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/vue/charts"
)

var NewUserValue *newUserValue

type newUserValue struct {
	*charts.Value
	srv *services.UserService
}

func NewNewUserValue() *newUserValue {
	if NewUserValue == nil {
		NewUserValue = &newUserValue{
			Value: charts.NewValue(),
			srv:   services.MakeUserService(),
		}
	}
	return NewUserValue
}

func (v newUserValue) Columns() []string {
	return []string{}
}

func (v newUserValue) HttpHandle(ctx *gin.Context) (rows interface{}, err error) {
	count := v.srv.TodayNewUserCount(ctx)
	return count, nil
}

func (newUserValue) Name() string {
	return "当日新增用户"
}
