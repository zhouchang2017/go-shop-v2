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
	register(NewBrand)
}

type Brand struct {
	vue.AbstractResource
	model *models.Brand
	rep   *repositories.BrandRep
}

func NewBrand(model *models.Brand, rep *repositories.BrandRep) *Brand {
	return &Brand{model: model, rep: rep}
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

func (b *Brand) IndexQuery(ctx *gin.Context, request *request.IndexRequest) {
	request.SetSearchField("name")
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
	return "i-inbox"
}