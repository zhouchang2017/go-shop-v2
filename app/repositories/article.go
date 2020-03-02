package repositories

import (
	"go-shop-v2/pkg/repository"
)

type ArticleRep struct {
	repository.IRepository
}

func NewArticleRep(rep repository.IRepository) *ArticleRep {
	return &ArticleRep{rep}
}
