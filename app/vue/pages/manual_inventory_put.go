package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"net/http"
)

var ManualInventoryPut *manualInventoryPut

// 自定义页面，库存操作入库
type manualInventoryPut struct {
	router  contracts.Router
	service *services.ManualInventoryActionService
}

func NewManualInventoryPut() *manualInventoryPut {
	con := mongodb.GetConFn()
	if ManualInventoryPut == nil {
		ManualInventoryPut = &manualInventoryPut{
			service: services.NewManualInventoryActionService(
				repositories.NewManualInventoryActionRep(con),
				repositories.NewShopRep(con),
				repositories.NewItemRep(con),
			),
		}
	}
	return ManualInventoryPut
}

func (this manualInventoryPut) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (this manualInventoryPut) VueRouter() contracts.Router {
	if this.router == nil {
		router := core.NewRouter()
		router.RouterPath = "inventory_actions/new"
		router.Name = "inventory_actions.create"
		router.RouterComponent = "inventories/Create"
		router.Hidden = true
		router.WithMeta("ResourceName", "inventory_actions")
		router.WithMeta("Title", this.Title())
		router.WithMeta("PutEndpoint", this.putEndpoint())
		router.WithMeta("TakeEndpoint", this.takeEndpoint())
		inventory := models.Inventory{}
		router.WithMeta("InventoryStatus", inventory.StatusOkMap())
		this.router = router
	}
	return this.router
}

func (this manualInventoryPut) HttpHandles(router gin.IRouter) {
	// 入库处理
	router.POST(this.putEndpoint(), func(c *gin.Context) {
		form := &services.InventoryActionPutOption{}
		if err := c.ShouldBind(form); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		user := ctx.GetUser(c)
		if admin, ok := user.(*models.Admin); ok {
			inventoryAction, err := this.service.Put(c, form, admin)
			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"redirect": fmt.Sprintf("/inventory_actions/%s", inventoryAction.GetID()),
			})
			return
		}
		err2.ErrorEncoder(nil, err2.Err401, c.Writer)
		return
	})

	// 出库处理
	router.POST(this.takeEndpoint(), func(c *gin.Context) {
		form := &services.InventoryActionTakeOption{}
		if err := c.ShouldBind(form); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		user := ctx.GetUser(c)
		if admin, ok := user.(*models.Admin); ok {
			inventoryAction, err := this.service.Take(c, form, admin)
			if err != nil {
				err2.ErrorEncoder(nil, err, c.Writer)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"redirect": fmt.Sprintf("/inventory_actions/%s", inventoryAction.GetID()),
			})
			return
		}
		err2.ErrorEncoder(nil, err2.Err401, c.Writer)
		return
	})
}

func (this manualInventoryPut) putEndpoint() string {
	return "/inventory_actions/put"
}

func (this manualInventoryPut) takeEndpoint() string {
	return "/inventory_actions/take"
}

func (this manualInventoryPut) Title() string {
	return "库存操作"
}

func (this manualInventoryPut) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}
