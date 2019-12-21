package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue"
)

func init()  {
	register(NewInventoryResource)
}

// 库存管理
type Inventory struct {
	vue.AbstractResource
	model   *models.Inventory
	rep     *repositories.InventoryRep
	service *services.InventoryService
}

func NewInventoryResource(model *models.Inventory, rep *repositories.InventoryRep, service *services.InventoryService) *Inventory {
	return &Inventory{model: model, rep: rep, service: service}
}

func (i *Inventory) IndexQuery(ctx *gin.Context, request *request.IndexRequest) {

}

func (i *Inventory) Model() interface{} {
	return i.model
}

func (i *Inventory) Repository() repository.IRepository {
	return i.rep
}

func (i Inventory) Make(model interface{}) vue.Resource {
	return &Inventory{model: model.(*models.Inventory)}
}

func (i *Inventory) SetModel(model interface{}) {
	i.model = model.(*models.Inventory)
}

func (i Inventory) Title() string {
	return "库存管理"
}

func (this Inventory) Icon() string {
	return "i-box"
}

func (Inventory) Group() string {
	return "Shop"
}