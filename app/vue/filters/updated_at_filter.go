package filters

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/filters"
)

type UpdatedAtFilter struct {
	*filters.DateTimeFilter
}

func NewUpdatedAtFilter() *UpdatedAtFilter {
	return &UpdatedAtFilter{DateTimeFilter: filters.NewDateTimeFilter()}
}

func (this UpdatedAtFilter) Key() string {
	return "updated_at"
}

func (this UpdatedAtFilter) Name() string {
	return "创建日期"
}



func (this UpdatedAtFilter) Options(ctx *gin.Context) []contracts.FilterOption {
	return []contracts.FilterOption{}
}
