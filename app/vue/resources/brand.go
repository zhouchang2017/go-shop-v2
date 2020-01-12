package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
)

type Brand struct {
	core.AbstractResource
	model   interface{}
	service *services.BrandService
}

func (b *Brand) Destroy(ctx *gin.Context, id string) (err error) {
	return b.service.Delete(ctx, id)
}

func (b *Brand) Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error) {
	brand, err := b.service.Update(ctx, model.(*models.Brand), data["name"].(string))
	if err != nil {
		return "", err
	}
	return core.CreatedRedirect(b, brand.GetID()), nil
}

func (b *Brand) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	brand, err := b.service.Create(ctx, data["name"].(string))
	if err != nil {
		return "", err
	}
	return core.CreatedRedirect(b, brand.GetID()), nil
}

func (b *Brand) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return b.service.FindById(ctx, id)
}

func (b *Brand) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return b.service.Pagination(ctx, req)
}

func (b *Brand) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}

func (b *Brand) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (b *Brand) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (b *Brand) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (b *Brand) Policy() interface{} {
	return nil
}

func (b *Brand) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("名称", "Name", fields.SetRules([]*fields.FieldRule{
				{Rule: "required"},
			})),
			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),
		}
	}
}

func NewBrandResource() *Brand {
	return &Brand{model: &models.Brand{}, service: services.MakeBrandService()}
}

func (b *Brand) Model() interface{} {
	return b.model
}

func (b Brand) Make(model interface{}) contracts.Resource {
	return &Brand{
		service: b.service,
		model:   model,
	}
}

func (b *Brand) SetModel(model interface{}) {
	b.model = model.(*models.Brand)
}

func (b Brand) Title() string {
	return "品牌"
}

func (b Brand) Group() string {
	return "Product"
}

func (b Brand) Icon() string {
	return "icons-inbox"
}
