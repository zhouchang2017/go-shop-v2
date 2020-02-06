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

type Topic struct {
	core.AbstractResource
	model   interface{}
	service *services.TopicService
}

func (t *Topic) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return t.service.FindById(ctx, id)
}

func (t *Topic) Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error) {
	option := services.TopicOption{}
	if err := mapstructure.Decode(data, &option); err != nil {
		return "", err
	}

	topic, err := t.service.Update(ctx, model.(*models.Topic), option)

	if err != nil {
		return "", err
	}

	return core.UpdatedRedirect(t, topic.GetID()), nil
}

func (t *Topic) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	option := services.TopicOption{}
	if err := mapstructure.Decode(data, &option); err != nil {
		return "", err
	}

	topic, err := t.service.Create(ctx, option)

	return core.CreatedRedirect(t, topic.GetID()), nil
}

func (t *Topic) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return t.service.Pagination(ctx, req)
}

func (t Topic) Title() string {
	return "话题"
}

func (t Topic) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("标题", "Title"),
			fields.NewTextField("副标题", "ShortTitle").Textarea(),
			fields.NewAvatar("封面", "Avatar").Rounded().URL(),
			fields.NewRichTextField("正文", "Content").UseQiniu(),
			fields.NewTextField("权重", "Sort").Min(0).Max(9999).InputNumber(),
			fields.NewRelationsField(&Product{}, "ProductIds").WithName("关联产品").Multiple().Searchable(),
			fields.NewDateTime("更新时间", "UpdatedAt"),
		}
	}
}

func (t Topic) Model() interface{} {
	return t.model
}

func (t Topic) Make(mode interface{}) contracts.Resource {
	return &Topic{
		model:   mode,
		service: t.service,
	}
}

func (t *Topic) SetModel(model interface{}) {
	t.model = model
}

func NewTopicResource() *Topic {
	return &Topic{model: &models.Topic{}, service: services.MakeTopicService()}
}
