package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
)

func init() {
	register(NewProductResource)
}

type Product struct {
	model   interface{}
	rep     *repositories.ProductRep
	service *services.ProductService
}

func (this *Product) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}

func (this *Product) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (this *Product) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (this *Product) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (this *Product) Policy() interface{} {
	return nil
}

func (this *Product) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{}
	}
}

func NewProductResource(rep *repositories.ProductRep, service *services.ProductService) *Product {
	return &Product{model: &models.Product{}, rep: rep, service: service}
}

// 自定义详情页数据
func (this *Product) CustomResourceHttpShow(ctx *gin.Context, id string) (model interface{}, err error) {
	return this.service.FindByIdWithItems(ctx, id)
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

func (this *Product) UpdateFormParse(ctx *gin.Context, model interface{}) (entity interface{}, err error) {
	form := &ProductForm{}
	err = ctx.ShouldBind(form)
	if err != nil {
		return nil, err
	}
	product := model.(*models.Product)

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
	return product, nil
}

func (this *Product) CreateFormParse(ctx *gin.Context) (entity interface{}, err error) {
	form := &ProductForm{}
	err = ctx.ShouldBind(form)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		Name:         form.Name,
		Code:         form.Code,
		Brand:        form.Brand,
		Category:     form.Category,
		Options:      form.Options,
		Attributes:   form.Attributes,
		Description:  form.Description,
		Price:        form.Price,
		FakeSalesQty: form.FakeSalesQty,
		Images:       form.Images,
		OnSale:       form.OnSale,
		Items:        form.Items,
	}

	return product, nil
}

func (this *Product) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	return nil
}

func (this *Product) Model() interface{} {
	return this.model
}


func (this Product) Make(model interface{}) contracts.Resource {
	return &Product{
		rep:     this.rep,
		service: this.service,
		model:   model,
	}
}

func (this *Product) SetModel(model interface{}) {
	this.model = model.(*models.Product)
}

func (this Product) Title() string {
	return "产品"
}

func (this Product) Group() string {
	return "Product"
}

func (this Product) Icon() string {
	return "icons-box"
}
