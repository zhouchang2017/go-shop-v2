package pages

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"net/http"
)

var ProductCreatePage *productCreatePage

type productCreatePage struct {
	brandRep    *repositories.BrandRep
	categoryRep *repositories.CategoryRep
	productService  *services.ProductService
}

func NewProductCreatePage() *productCreatePage {
	if ProductCreatePage == nil {
		con := mongodb.GetConFn()
		ProductCreatePage = &productCreatePage{
			brandRep:    repositories.NewBrandRep(con),
			categoryRep: repositories.NewCategoryRep(con),
			productService:  services.MakeProductService(),
		}
	}
	return ProductCreatePage
}

func (this *productCreatePage) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (this *productCreatePage) VueRouter() contracts.Router {
	router := core.NewRouter()
	router.RouterPath = "products/new"
	router.Name = "products.create"
	router.RouterComponent = "products/Create"
	router.Hidden = true
	router.WithMeta("ResourceName", "products")
	router.WithMeta("Title", this.Title())
	return router
}

func (this *productCreatePage) getBrands(ctx context.Context) ([]*models.Brand, error) {
	all := <-this.brandRep.FindAll(ctx)
	if all.Error != nil {
		return nil, all.Error
	}
	return all.Result.([]*models.Brand), nil
}

func (this *productCreatePage) getCategories(ctx context.Context) ([]*models.Category, error) {
	all := <-this.categoryRep.FindAll(ctx)
	if all.Error != nil {
		return nil, all.Error
	}
	return all.Result.([]*models.Category), nil
}

func (this *productCreatePage) HttpHandles(router gin.IRouter) {
	// 创建产品关联数据
	router.GET("creation-info/products", func(c *gin.Context) {
		if !this.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			c.AbortWithStatus(403)
			return
		}
		brands, err := this.getBrands(c)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		categories, err := this.getCategories(c)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"brands":     brands,
			"categories": categories,
		})
	})
	// 创建产品api
	router.POST("products", func(c *gin.Context) {
		if !this.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			c.AbortWithStatus(403)
			return
		}

		option := services.ProductCreateOption{}
		err := c.ShouldBind(&option)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		product, err := this.productService.Create(c, option)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"redirect": fmt.Sprintf("/products/%s", product.GetID())})
	})
}

type ProductForm struct {
	Name         string                     `json:"name" form:"name" binding:"required,max=255"`
	Code         string                     `json:"code" form:"code" binding:"required,max=255"`
	Brand        *models.AssociatedBrand    `json:"brand" form:"brand"`
	Category     *models.AssociatedCategory `json:"category" form:"category"`
	Attributes   []*models.ProductAttribute `json:"attributes" form:"attributes"`
	Options      []*models.ProductOption    `json:"options" form:"options"`
	Items        []*models.Item             `json:"items"`
	Description  string                     `json:"description"`
	Price        int64                      `json:"price"`
	FakeSalesQty int64                      `json:"fake_sales_qty" form:"fake_sales_qty"`
	Images       []*qiniu.Resource          `json:"images" form:"images"`
	OnSale       bool                       `json:"on_sale" form:"on_sale"`
}

func (this *productCreatePage) Title() string {
	return "创建产品"
}

func (this *productCreatePage) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}
