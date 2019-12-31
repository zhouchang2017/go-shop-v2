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

var ProductUpdatePage *productUpdatePage

type productUpdatePage struct {
	productService *services.ProductService
	productRep  *repositories.ProductRep
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
		if err!=nil {
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

		form := &ProductForm{}
		err := c.ShouldBind(form)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		result := <-this.productRep.FindById(c, c.Param("Product"))
		if result.Error != nil {
			err2.ErrorEncoder(nil, result.Error, c.Writer)
			return
		}

		product := result.Result.(*models.Product)

		product.Name = form.Name
		product.Brand = form.Brand
		product.Items = form.Items
		product.Options = form.Options
		product.Attributes = form.Attributes
		product.Description = form.Description
		product.Price = form.Price
		product.FakeSalesQty = form.FakeSalesQty
		product.Images = form.Images
		product.OnSale = form.OnSale

		updated := <-this.productRep.Save(c, product)
		if updated.Error != nil {
			err2.ErrorEncoder(nil, updated.Error, c.Writer)
			return
		}

		c.JSON(http.StatusOK, gin.H{"redirect": fmt.Sprintf("/products/%s", product.GetID())})
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
		con := mongodb.GetConFn()
		itemRep := repositories.NewItemRep(con)
		productRep:=  repositories.NewProductRep(con, itemRep)

		ProductUpdatePage = &productUpdatePage{
			productService:services.NewProductService(productRep),
			productRep:  productRep,
		}
	}
	return ProductUpdatePage
}
