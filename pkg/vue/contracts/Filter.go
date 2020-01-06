package contracts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/request"
)

type Filter interface {
	Apply(ctx *gin.Context, value interface{}, request *request.IndexRequest) error
	Key() string
	Name() string
	DefaultValue(ctx *gin.Context) interface{}
	Options(ctx *gin.Context) []FilterOption
	Element
}

type FilterOption interface {
	Label() string
	Value() interface{}
}
