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

type ArticleService struct {
	rep *repositories.ArticleRep
}

func NewArticleService(rep *repositories.ArticleRep) *ArticleService {
	return &ArticleService{rep: rep}
}

// 列表
func (this *ArticleService) Pagination(ctx context.Context, req *request.IndexRequest) (articles []*models.Article, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Article), results.Pagination, nil
}

// 简单列表
func (this *ArticleService) SimplePagination(ctx context.Context, page int64, perPage int64) (articles []*models.Article, pagination response.Pagination, err error) {
	req:=&request.IndexRequest{
		Page:           page,
		PerPage:        perPage,
		OrderBy:        "sort",
		OrderDirection: -1,
	}
	// 不展现 content,product_id
	req.Hidden = "content,product_id"
	// 只搜索第一张图片
	req.AppendProjection("photos", bson.M{"$slice": 1})
	results := <-this.rep.Pagination(ctx,req )
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Article), results.Pagination, nil
}

// 总计数量
func (this *ArticleService) Count(ctx context.Context) int64 {
	count := <-this.rep.Count(ctx, bson.M{})
	if count.Error != nil {
		return 0
	}
	return count.Result
}

// 表单结构
type ArticleOption struct {
	Title      string   `json:"title"`
	ShortTitle string   `json:"short_title" mapstructure:"short_title"`
	Photos     []qiniu.Image `json:"photos"`
	Content    string   `json:"content"`
	ProductId  string   `json:"product_id" mapstructure:"product_id"`
	Sort       int64    `json:"sort"`
}

// 创建文章
func (this *ArticleService) Create(ctx context.Context, opt ArticleOption) (article *models.Article, err error) {
	created := <-this.rep.Create(ctx, &models.Article{
		Title:      opt.Title,
		ShortTitle: opt.ShortTitle,
		Photos:     opt.Photos,
		Content:    opt.Content,
		ProductId:  opt.ProductId,
		Sort:       opt.Sort,
	})

	if created.Error != nil {
		return nil, created.Error
	}

	return created.Result.(*models.Article), nil
}

// 更新文章
func (this *ArticleService) Update(ctx context.Context, model *models.Article, opt ArticleOption) (article *models.Article, err error) {
	model.Title = opt.Title
	model.ShortTitle = opt.ShortTitle
	model.Photos = opt.Photos
	model.Content = opt.Content
	model.ProductId = opt.ProductId
	model.Sort = opt.Sort
	saved := <-this.rep.Save(ctx, model)
	if saved.Error != nil {
		return nil, saved.Error
	}
	return saved.Result.(*models.Article), nil
}

// 文章详情
func (this *ArticleService) FindById(ctx context.Context, id string) (article *models.Article, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		return nil, byId.Error
	}
	return byId.Result.(*models.Article), nil
}

// 删除
func (this *ArticleService) Delete(ctx context.Context, id string) (err error) {
	return <-this.rep.Delete(ctx, id)
}

// 还原
func (this *ArticleService) Restore(ctx context.Context, id string) (article *models.Article, err error) {
	restored := <-this.rep.Restore(ctx, id)
	if restored.Error != nil {
		return nil, restored.Error
	}
	return restored.Result.(*models.Article), nil
}
