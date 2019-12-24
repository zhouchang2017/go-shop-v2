package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue"
)

func init() {
	register(NewItemResource)
}

type Item struct {
	vue.AbstractResource
	model *models.Item
	rep   *repositories.ItemRep
}

func (i *Item) Model() interface{} {
	return i.model
}

func (i *Item) Repository() repository.IRepository {
	return i.rep
}

func (i Item) Make(model interface{}) vue.Resource {
	return &Item{model: model.(*models.Item)}
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

func (Item) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	request.SetSearchField("code")
	return nil
}

func NewItemResource(model *models.Item, rep *repositories.ItemRep) *Item {
	return &Item{model: model, rep: rep}
}

func (i *Item)DisplayInNavigation(ctx *gin.Context) bool  {
	return false
}