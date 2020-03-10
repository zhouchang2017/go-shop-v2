package pages

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
)

var PromotionCreatePage *promotionCreatePage

// 促销计划自定义创建页面
type promotionCreatePage struct {
	promotionService *services.PromotionService
	productService   *services.ProductService
}

func NewPromotionCreatePage() *promotionCreatePage {
	if PromotionCreatePage == nil {
		PromotionCreatePage = &promotionCreatePage{promotionService: services.MakePromotionService(), productService: services.MakeProductService()}
	}
	return PromotionCreatePage
}

func (p *promotionCreatePage) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (p *promotionCreatePage) VueRouter() contracts.Router {
	router := core.NewRouter()
	router.RouterPath = "promotions/new"
	router.Name = "promotions.create"
	router.RouterComponent = "promotions/Create"
	router.Hidden = true
	router.WithMeta("ResourceName", "promotions")
	router.WithMeta("Title", p.Title())
	return router
}

func (p promotionCreatePage) HttpHandles(router gin.IRouter) {

}

func (p promotionCreatePage) Title() string {
	return "新建促销计划"
}

func (p promotionCreatePage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}
