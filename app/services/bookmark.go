package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
)

type BookmarkService struct {
	rep *repositories.BookmarkRep
}

func (this *BookmarkService) Pagination(ctx context.Context, req *request.IndexRequest) (bookmarks []*models.Bookmark, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	bookmarks = results.Result.([]*models.Bookmark)
	if len(bookmarks) == 0 {
		bookmarks = []*models.Bookmark{}
	}
	pagination = results.Pagination
	return
}

// 添加到收藏夹
func (this *BookmarkService) Add(ctx context.Context, product *models.Product, userId string) (bookmark *models.Bookmark, err error) {
	// 检测是否已在购物车
	hasOne := <-this.rep.FindOne(ctx, bson.M{
		"user_id":    userId,
		"product.id": product.GetID(),
	})
	if hasOne.Error != nil {
		// 不存在，创建
		model := &models.Bookmark{
			UserId:  userId,
			Product: product.ToAssociated(),
			Enabled: true,
		}
		created := <-this.rep.Create(ctx, model)
		if created.Error != nil {
			// 创建失败
			return nil, created.Error
		}
		return created.Result.(*models.Bookmark), nil
	}
	return nil, err2.New(1002, "已存在")
}

// 从收藏夹移除
func (this *BookmarkService) Delete(ctx context.Context, ids ...string) error {
	ctx = ctx2.WithForce(ctx, true)
	err := <-this.rep.DeleteMany(ctx, ids...)
	return err
}

func NewBookmarkService(rep *repositories.BookmarkRep) *BookmarkService {
	return &BookmarkService{rep: rep}
}
