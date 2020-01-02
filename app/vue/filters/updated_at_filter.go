package filters

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/filters"
)

type UpdatedAtFilter struct {
	*filters.DateTimeFilter
}



func NewUpdatedAtFilter() *UpdatedAtFilter {
	return &UpdatedAtFilter{DateTimeFilter: filters.NewDateTimeFilter()}
}

func (this UpdatedAtFilter) Apply(ctx *gin.Context, value interface{}, request *request.IndexRequest) error {
	spew.Dump(value)
	return nil
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
