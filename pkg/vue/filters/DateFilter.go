package filters

import (
	"github.com/gin-gonic/gin"
	"time"
)

type DateFilter struct {
	*AbstractSelectFilter
}

func NewDateFilter() *DateFilter {
	return &DateFilter{AbstractSelectFilter: NewAbstractSelectFilter()}
}

func (this DateFilter) Component() string {
	return "date-filter"
}

func (this DateFilter) DefaultValue(ctx *gin.Context) interface{} {
	return time.Now().Format("2006-01-02 15:04:05")
}