package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/fields"
)

type PromotionItemResource struct {
	service *services.PromotionItemService
	model   interface{}
}

func NewPromotionItemResource() *PromotionItemResource {
	return &PromotionItemResource{
		service: services.MakePromotionItemService(),
		model:   &models.PromotionItem{},
	}
}

func (this *PromotionItemResource) Pagination(ctx *gin.Context,
	req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	req.SetSearchField("product_id")
	promotionId := ctx.Query("promotion_id")
	if promotionId != "" {
		req.AppendFilter("promotion.id", promotionId)
	}
	return this.service.Pagination(ctx, req)
}

func (PromotionItemResource) Title() string {
	return "促销商品"
}

func (PromotionItemResource) Group() string {
	return "Product"
}

func (*PromotionItemResource) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return false
}

func (*PromotionItemResource) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (*PromotionItemResource) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (*PromotionItemResource) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return false
}

func (*PromotionItemResource) Policy() interface{} {
	return nil
}

func (*PromotionItemResource) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewAvatar("缩略图", "Product.Avatar").Rounded(),
			fields.NewTextField("产品名称", "Product.Name"),
			fields.NewTextField("货号", "Product.Code"),
			fields.NewCurrencyField("吊牌价", "Product.Price"),
			fields.NewTable("SKU", "Units", func() []contracts.Field {
				return []contracts.Field{
					fields.NewAvatar("缩略图", "Item.Avatar").Rounded(),
					fields.NewTextField("编码", "Item.Code"),
					fields.NewLabelsFields("销售属性", "Item.OptionValues").Label("name"),
					fields.NewCurrencyField("吊牌加", "Item.Price"),
					fields.NewCurrencyField("折后价", "Price"),
				}
			}, fields.SetExpand(true)),
		}
	}
}

func (this *PromotionItemResource) Model() interface{} {
	return this.model
}

func (this *PromotionItemResource) Make(mode interface{}) contracts.Resource {
	return &PromotionItemResource{
		model:   mode,
		service: this.service,
	}
}

func (this *PromotionItemResource) SetModel(model interface{}) {
	this.model = model
}

func (*PromotionItemResource) Lenses() []contracts.Lens {
	return []contracts.Lens{}
}

func (*PromotionItemResource) Pages() []contracts.Page {
	return []contracts.Page{}
}

func (*PromotionItemResource) Filters(ctx *gin.Context) []contracts.Filter {
	return []contracts.Filter{}
}

func (*PromotionItemResource) Actions(ctx *gin.Context) []contracts.Action {
	return []contracts.Action{}
}

func (*PromotionItemResource) Cards(ctx *gin.Context) []contracts.Card {
	return []contracts.Card{}
}
