package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
)

type Article struct {
	core.AbstractResource
	model   interface{}
	service *services.ArticleService
}

func NewArticleResource() *Article {
	return &Article{model: &models.Article{}, service: services.MakeArticleService()}
}

func (a *Article) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	option := services.ArticleOption{}
	if err := mapstructure.Decode(data, &option); err != nil {
		return "", err
	}

	article, err := a.service.Create(ctx, option)

	return core.CreatedRedirect(a, article.GetID()), nil
}

func (a *Article) Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error) {
	option := services.ArticleOption{}
	if err := mapstructure.Decode(data, &option); err != nil {
		return "", err
	}
	article := model.(*models.Article)
	article2, err := a.service.Update(ctx, article, option)

	if err != nil {
		return "", err
	}

	return core.UpdatedRedirect(a, article2.GetID()), nil
}

func (a *Article) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return a.service.Pagination(ctx, req)
}

func (a *Article) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return a.service.FindById(ctx, id)
}

func (a Article) Title() string {
	return "文章"
}

func (a Article) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("标题", "Title"),
			fields.NewTextField("副标题", "ShortTitle").Textarea(),
			fields.NewImageField("图集", "Photos").Multiple().Limit(5).Rounded().URL(),
			fields.NewRichTextField("正文", "Content").UseQiniu(),
			fields.NewTextField("权重", "Sort").Min(0).Max(9999).InputNumber(),
			fields.NewRelationsField(&Product{}, "ProductId").WithName("关联产品").Searchable(),
			fields.NewDateTime("更新时间", "UpdatedAt"),
		}
	}
}

func (a Article) Model() interface{} {
	return a.model
}

func (a Article) Make(mode interface{}) contracts.Resource {
	return &Article{
		model:   mode,
		service: a.service,
	}
}

func (a *Article) SetModel(model interface{}) {
	a.model = model
}

func (order *Article) Icon() string {
	return "icons-exclamation"
}