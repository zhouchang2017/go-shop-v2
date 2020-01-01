package contracts

import "github.com/gin-gonic/gin"

type Filter interface {
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