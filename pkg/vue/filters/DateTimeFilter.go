package filters

import (
	"github.com/gin-gonic/gin"
)

type DateTimeFilter struct {
	*AbstractSelectFilter
}

func NewDateTimeFilter() *DateTimeFilter {
	return &DateTimeFilter{AbstractSelectFilter: NewAbstractSelectFilter()}
}

func (this DateTimeFilter) Component() string {
	return "date-time-filter"
}

func (this DateTimeFilter) DefaultValue(ctx *gin.Context) interface{} {
	return []interface{}{}
}