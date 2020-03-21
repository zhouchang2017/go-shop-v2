package repositories

import "go-shop-v2/pkg/repository"

type CommentRep struct {
	repository.IRepository
}

func NewCommentRep(rep repository.IRepository) *CommentRep {
	return &CommentRep{rep}
}
