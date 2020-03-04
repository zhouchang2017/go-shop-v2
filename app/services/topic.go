package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
)

type TopicService struct {
	rep *repositories.TopicRep
}

func NewTopicService(rep *repositories.TopicRep) *TopicService {
	return &TopicService{rep: rep}
}

// 列表
func (this *TopicService) Pagination(ctx context.Context, req *request.IndexRequest) (articles []*models.Topic, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Topic), results.Pagination, nil
}

// 总计数量
func (this *TopicService) Count(ctx context.Context) int64 {
	count := <-this.rep.Count(ctx, bson.M{})
	if count.Error != nil {
		return 0
	}
	return count.Result
}

// 简单列表
func (this *TopicService) SimplePagination(ctx context.Context, page int64, perPage int64) (articles []*models.Topic, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, &request.IndexRequest{
		Page:           page,
		PerPage:        perPage,
		OrderBy:        "sort",
		OrderDirection: -1,
	})
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Topic), results.Pagination, nil
}

// 表单结构
type TopicOption struct {
	Title      string   `json:"title"`
	ShortTitle string   `json:"short_title" mapstructure:"short_title"`
	Avatar     qiniu.Image   `json:"avatar"`
	Content    string   `json:"content"`
	ProductIds []string `json:"product_ids" mapstructure:"product_ids"`
	Sort       int64    `json:"sort"`
}

// 创建话题
func (this *TopicService) Create(ctx context.Context, opt TopicOption) (topic *models.Topic, err error) {
	created := <-this.rep.Create(ctx, &models.Topic{
		Title:      opt.Title,
		ShortTitle: opt.ShortTitle,
		Avatar:     opt.Avatar,
		Content:    opt.Content,
		ProductIds: opt.ProductIds,
		Sort:       opt.Sort,
	})

	if created.Error != nil {
		return nil, created.Error
	}

	return created.Result.(*models.Topic), nil
}

// 更新话题
func (this *TopicService) Update(ctx context.Context, model *models.Topic, opt TopicOption) (topic *models.Topic, err error) {
	model.Title = opt.Title
	model.ShortTitle = opt.ShortTitle
	model.Avatar = opt.Avatar
	model.Content = opt.Content
	model.ProductIds = opt.ProductIds
	model.Sort = opt.Sort
	saved := <-this.rep.Save(ctx, model)
	if saved.Error != nil {
		return nil, saved.Error
	}
	return saved.Result.(*models.Topic), nil
}

// 话题详情
func (this *TopicService) FindById(ctx context.Context, id string) (topic *models.Topic, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		return nil, byId.Error
	}
	return byId.Result.(*models.Topic), nil
}

// 删除
func (this *TopicService) Delete(ctx context.Context, id string) (err error) {
	return <-this.rep.Delete(ctx, id)
}

// 还原
func (this *TopicService) Restore(ctx context.Context, id string) (article *models.Topic, err error) {
	restored := <-this.rep.Restore(ctx, id)
	if restored.Error != nil {
		return nil, restored.Error
	}
	return restored.Result.(*models.Topic), nil
}
