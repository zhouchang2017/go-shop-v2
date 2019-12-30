package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
)

func init() {
	register(NewBrandResource)
}

type Brand struct {
	model interface{}
	rep   *repositories.BrandRep
}

func (b *Brand) Destroy(ctx *gin.Context, id string) (err error) {
	return <-b.rep.Delete(ctx, id)
}

type brandForm struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func (b *Brand) Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error) {
	form := &brandForm{}
	if err := mapstructure.Decode(data, form); err != nil {
		return "", err
	}
	brand := model.(*models.Brand)
	brand.Name = form.Name
	saved := <-b.rep.Save(ctx, brand)
	if saved.Error != nil {
		return "", saved.Error
	}
	return core.CreatedRedirect(b, brand.GetID()), nil
}

func (b *Brand) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	form := &brandForm{}
	if err := mapstructure.Decode(data, form); err != nil {
		return "", err
	}
	brand := &models.Brand{Name: form.Name}
	created := <-b.rep.Create(ctx, brand)
	if created.Error != nil {
		return "", created.Error
	}
	return core.CreatedRedirect(b, created.Id), nil
}

func (b *Brand) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	result := <-b.rep.FindById(ctx, id)
	return result.Result, result.Error
}

func (b *Brand) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	results := <-b.rep.Pagination(ctx, req)
	return results.Result, results.Pagination, results.Error
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

func NewBrandResource(rep *repositories.BrandRep) *Brand {
	return &Brand{model: &models.Brand{}, rep: rep}
}


func (b *Brand) Model() interface{} {
	return b.model
}

func (b Brand) Make(model interface{}) contracts.Resource {
	return &Brand{
		rep:   b.rep,
		model: model,
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
