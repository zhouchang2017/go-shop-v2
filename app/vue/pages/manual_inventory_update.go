package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"net/http"
)

var ManualInventoryUpdatePage *manualInventoryUpdatePage
// 自定义页面，更新库存操作
type manualInventoryUpdatePage struct {
	router  contracts.Router
	service *services.ManualInventoryActionService
}

func (this manualInventoryUpdatePage) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (this manualInventoryUpdatePage) VueRouter() contracts.Router {
	if this.router == nil {
		router := core.NewRouter()
		router.RouterPath = "inventory_actions/:id/edit"
		router.Name = "inventory_actions.edit"
		router.RouterComponent = "inventories/Edit"
		router.Hidden = true
		router.WithMeta("ResourceName", "inventory_actions")
		router.WithMeta("Title", this.Title())
		inventory := models.Inventory{}
		router.WithMeta("InventoryStatus", inventory.StatusOkMap())
		this.router = router
	}
	return this.router
}

func (this manualInventoryUpdatePage) HttpHandles(router gin.IRouter) {
	router.GET("api/inventory_actions/:InventoryAction/inventories", func(c *gin.Context) {
		action, err := this.service.FindByIdWithInventory(c, c.Param("InventoryAction"))
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		c.JSON(http.StatusOK, action)
		return
	})
	// 入库更新处理
	router.PUT("inventory_actions/:InventoryAction/put", func(c *gin.Context) {
		form := &services.InventoryActionPutOption{}
		if err := c.ShouldBind(form); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		user := ctx.GetUser(c)
		if admin, ok := user.(*models.Admin); ok {
			inventoryAction, err := this.service.UpdatePut(c, c.Param("InventoryAction"), form, admin)
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
	// 出库更新处理
	router.PUT("inventory_actions/:InventoryAction/take", func(c *gin.Context) {
		form := &services.InventoryActionTakeOption{}
		if err := c.ShouldBind(form); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		user := ctx.GetUser(c)
		if admin, ok := user.(*models.Admin); ok {
			inventoryAction, err := this.service.UpdateTake(c, c.Param("InventoryAction"), form, admin)
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

func (this manualInventoryUpdatePage) Title() string {
	return "更新库存操作"
}

func (this manualInventoryUpdatePage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}

func NewManualInventoryUpdatePage() *manualInventoryUpdatePage {
	if ManualInventoryUpdatePage == nil {
		ManualInventoryUpdatePage = &manualInventoryUpdatePage{
			service: services.MakeManualInventoryActionService(),
		}
	}
	return ManualInventoryUpdatePage
}
