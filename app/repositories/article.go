package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

type ArticleRep struct {
	*mongoRep
}

func NewArticleRep(con *mongodb.Connection) *ArticleRep {
	return &ArticleRep{NewBasicMongoRepositoryByDefault(&models.Article{}, con)}
}
