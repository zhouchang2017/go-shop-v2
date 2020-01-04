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

var ProductUpdatePage *productUpdatePage

type productUpdatePage struct {
	productService *services.ProductService
}

func (this *productUpdatePage) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (this *productUpdatePage) VueRouter() contracts.Router {
	router := core.NewRouter()
	router.RouterPath = "products/:id/edit"
	router.Name = "products.edit"
	router.RouterComponent = "products/Edit"
	router.Hidden = true
	router.WithMeta("ResourceName", "products")
	router.WithMeta("Title", this.Title())
	return router
}

func (this *productUpdatePage) HttpHandles(router gin.IRouter) {
	// 产品数据
	router.GET("products/:Product/api", func(c *gin.Context) {
		if !this.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			c.AbortWithStatus(403)
			return
		}
		product, err := this.productService.FindByIdWithItems(c, c.Param("Product"))
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		c.JSON(http.StatusOK, product)
	})
	// 更新产品api
	router.PUT("products/:Product", func(c *gin.Context) {
		if !this.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			c.AbortWithStatus(403)
			return
		}

		product, err := this.productService.FindById(c, c.Param("Product"))
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		form := services.ProductCreateOption{}

		if err := c.ShouldBind(&form); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		updatedProduct, err := this.productService.Update(c, product, form)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		c.JSON(http.StatusOK, gin.H{"redirect": fmt.Sprintf("/products/%s", updatedProduct.GetID())})
	})
}

func (this *productUpdatePage) Title() string {
	return "编辑产品"
}

func (this *productUpdatePage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}

func NewProductUpdatePage() *productUpdatePage {
	if ProductUpdatePage == nil {
		ProductUpdatePage = &productUpdatePage{
			productService: services.MakeProductService(),
		}
	}
	return ProductUpdatePage
}
