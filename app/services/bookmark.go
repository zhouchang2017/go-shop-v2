package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
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

// 收藏夹产品列表
func (this *BookmarkService) Index(ctx context.Context, userId string, page int64, perPage int64) (ids []string, pagination response.Pagination, err error) {
	return this.rep.Index(ctx, userId, page, perPage)
}

// 添加到收藏夹
func (this *BookmarkService) Add(ctx context.Context, userId string, productId string) (err error) {
	return this.rep.Add(ctx, userId, productId)
}

// 从收藏夹移除
func (this *BookmarkService) Remove(ctx context.Context, userId string, ids ...string) error {
	return this.rep.Remove(ctx, userId, ids...)
}

// 收藏夹总数
func (this *BookmarkService) Count(ctx context.Context, userId string) (count int64) {
	return this.rep.Count(ctx, userId)
}

func (this *BookmarkService) FindByProductId(ctx context.Context, userId string, productId string) (bookmark *models.Bookmark) {
	return this.rep.FindByProductId(ctx, userId, productId)
}

func NewBookmarkService(rep *repositories.BookmarkRep) *BookmarkService {
	return &BookmarkService{rep: rep}
}
