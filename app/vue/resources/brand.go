package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue"
)

func init() {
	register(NewBrandResource)
}

type Brand struct {
	vue.AbstractResource
	model *models.Brand
	rep   *repositories.BrandRep
}

func NewBrandResource(rep *repositories.BrandRep) *Brand {
	return &Brand{model: &models.Brand{}, rep: rep}
}

// 列表页&详情页展示字段设置
func (s *Brand) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			vue.NewIDField(),
			vue.NewTextField("名称", "Name"),
			vue.NewDateTime("创建时间", "CreatedAt"),
			vue.NewDateTime("更新时间", "UpdatedAt"),

		}
	}
}

type brandForm struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func (b *Brand) UpdateFormParse(ctx *gin.Context, model interface{}) (entity interface{}, err error) {
	form := &brandForm{}
	err = ctx.ShouldBind(form)
	if err != nil {
		return nil, err
	}
	brand := model.(*models.Brand)
	brand.Name = form.Name
	return brand, nil
}

func (b *Brand) CreateFormParse(ctx *gin.Context) (entity interface{}, err error) {
	form := &brandForm{}
	err = ctx.ShouldBind(form)
	if err != nil {
		return nil, err
	}
	return &models.Brand{Name: form.Name}, nil
}

func (b *Brand) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	request.SetSearchField("name")
	return nil
}

func (b *Brand) Model() interface{} {
	return b.model
}

func (b *Brand) Repository() repository.IRepository {
	return b.rep
}

func (b Brand) Make(model interface{}) vue.Resource {
	return &Brand{model: model.(*models.Brand)}
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