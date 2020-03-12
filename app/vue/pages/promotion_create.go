package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"net/http"
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

func (p *promotionCreatePage) HttpHandles(router gin.IRouter) {
	// 创建促销api
	router.POST("promotions", func(c *gin.Context) {
		if !p.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			c.AbortWithStatus(403)
			return
		}

		option := services.PromotionCreateOption{}
		err := c.ShouldBind(&option)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		promotion, err := p.promotionService.Create(c, &option)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"redirect": fmt.Sprintf("/promotions/%s", promotion.GetID())})
	})
}

func (p promotionCreatePage) Title() string {
	return "新建促销计划"
}

func (p promotionCreatePage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}
