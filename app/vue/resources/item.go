package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
)

type Item struct {
	core.AbstractResource
	model interface{}
	rep   *repositories.ItemRep
}

func (i *Item) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	results := <-i.rep.Pagination(ctx, req)
	return results.Result, results.Pagination, results.Error
}

func (i *Item) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}

func (i *Item) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (i *Item) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (i *Item) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (i *Item) Policy() interface{} {
	return nil
}

func (i *Item) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{}
	}
}

func (i *Item) Model() interface{} {
	return i.model
}

func (i *Item) Repository() repository.IRepository {
	return i.rep
}

func (i Item) Make(model interface{}) contracts.Resource {
	return &Item{
		rep:   i.rep,
		model: model,
	}
}

func (i *Item) SetModel(model interface{}) {
	i.model = model.(*models.Item)
}

func (i Item) Title() string {
	return "SKU"
}

func (Item) Group() string {
	return "Product"
}

func NewItemResource() *Item {
	return &Item{model: &models.Item{}, rep: repositories.NewItemRep(mongodb.GetConFn())}
}
