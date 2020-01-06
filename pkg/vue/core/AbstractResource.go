package core

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/vue/contracts"
)

type AbstractResource struct {
}

func (AbstractResource) Group() string {
	return "App"
}

func (AbstractResource) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}

func (AbstractResource) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (AbstractResource) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (AbstractResource) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (AbstractResource) Policy() interface{} {
	return nil
}

func (AbstractResource) Filters(ctx *gin.Context) []contracts.Filter {
	return []contracts.Filter{}
}

func (AbstractResource) Lenses() []contracts.Lens {
	return []contracts.Lens{}
}

func (AbstractResource) Pages() []contracts.Page {
	return []contracts.Page{}
}

func (AbstractResource) Actions(ctx *gin.Context) []contracts.Action {
	return []contracts.Action{}
}
