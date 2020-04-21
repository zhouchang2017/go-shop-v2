package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type CommentService struct {
	rep *repositories.CommentRep
}

// 列表
func (this *CommentService) Pagination(ctx context.Context, req *request.IndexRequest) (comments []*models.Comment, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Comment), results.Pagination, nil
}

func NewCommentService(rep *repositories.CommentRep) *CommentService {
	return &CommentService{rep: rep}
}
