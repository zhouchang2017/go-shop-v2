package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/app/vue/pages"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
)

type Product struct {
	core.AbstractResource
	model   interface{}
	service *services.ProductService
}

// 自定义更新页
func (this *Product) UpdateComponent() contracts.Page {
	return pages.NewProductUpdatePage()
}

// 自定义创建页
func (this *Product) CreationComponent() contracts.Page {
	return pages.NewProductCreatePage()
}

// 实现列表页api
func (this *Product) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return this.service.Pagination(ctx, req)
}

// 实现详情页api
func (this *Product) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return this.service.FindByIdWithItems(ctx, id)
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
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("货号", "Code"),
			fields.NewTextField("商品名称", "Name"),
			fields.NewTextField("品牌", "Brand.Name"),
			fields.NewTextField("类目", "Category.Name"),
			fields.NewTextField("价格", "Price"),
			fields.NewTextField("销量", "TotalSalesQty"),

			fields.NewImageField("图集","Images"),

			fields.NewTable("基本属性", "Attributes", func() []contracts.Field {
				return []contracts.Field{
					fields.NewTextField("属性名", "Name"),
					fields.NewTextField("属性值", "Value"),
				}
			}),

			fields.NewTable("销售属性", "Options", func() []contracts.Field {
				return []contracts.Field{
					fields.NewTextField("ID", "Id"),
					fields.NewTextField("属性名", "Name"),
					fields.NewTextField("权重", "Sort"),
					fields.NewLabelsFields("属性值", "Values").Label("value").Tooltip("code"),
				}
			}),

			fields.NewTable("SKU", "Items", func() []contracts.Field {
				return []contracts.Field{
					fields.NewIDField(),
					fields.NewTextField("编码", "Code"),
					fields.NewTextField("价格", "Price"),
					fields.NewLabelsFields("销售属性", "OptionValues").Label("value").Tooltip("code"),
					fields.NewTextField("销量", "SalesQty"),
				}
			}),
		}
	}
}

func NewProductResource() *Product {
	return &Product{model: &models.Product{}, service: services.MakeProductService()}
}

// 自定义详情页数据
func (this *Product) CustomResourceHttpShow(ctx *gin.Context, id string) (model interface{}, err error) {
	return this.service.FindByIdWithItems(ctx, id)
}

func (this *Product) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	return nil
}

func (this *Product) Model() interface{} {
	return this.model
}

func (this Product) Make(model interface{}) contracts.Resource {
	return &Product{
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
