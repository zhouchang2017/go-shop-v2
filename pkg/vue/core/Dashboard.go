package core

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/vue/contracts"
)

type dashboard struct {
	AbstractResource
}

func (d dashboard) Icon() string {
	return "icons-home"
}

func (d dashboard) Title() string {
	return "Dashboard"
}

func (d dashboard) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{}
	}
}

func (d dashboard) Model() interface{} {
	return nil
}

func (d dashboard) Make(mode interface{}) contracts.Resource {
	return &dashboard{}
}

func (d dashboard) SetModel(model interface{}) {

}

func (dashboard) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return false
}

func (dashboard) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return false
}
