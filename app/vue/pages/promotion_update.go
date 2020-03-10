package pages

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/app/usecases"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"net/http"
)

var PromotionUpdatePage *promotionUpdatePage

// 促销计划自定义更新页面
type promotionUpdatePage struct {
	promotionService *services.PromotionService
	productService   *services.ProductService
}

func NewPromotionUpdatePage() *promotionUpdatePage {
	if PromotionUpdatePage == nil {
		PromotionUpdatePage = &promotionUpdatePage{
			productService:   services.MakeProductService(),
			promotionService: services.MakePromotionService(),
		}
	}
	return PromotionUpdatePage
}

func (p *promotionUpdatePage) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (p *promotionUpdatePage) VueRouter() contracts.Router {
	router := core.NewRouter()
	router.RouterPath = "promotions/:id/edit"
	router.Name = "promotions.edit"
	router.RouterComponent = "promotions/Edit"
	router.Hidden = true
	router.WithMeta("ResourceName", "promotions")
	router.WithMeta("Title", p.Title())
	return router
}

func (p *promotionUpdatePage) HttpHandles(router gin.IRouter) {
	// 促销计划数据
	router.GET("promotions/:Promotion/api", func(c *gin.Context) {
		if !p.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			c.AbortWithStatus(403)
			return
		}
		id := c.Param("Promotion")
		promotion, err := usecases.PromotionWithItemsAndProduct(c, id, p.promotionService, p.productService)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		c.JSON(http.StatusOK, promotion)
	})
}

func (p promotionUpdatePage) Title() string {
	return "更新促销计划"
}

func (p *promotionUpdatePage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}
