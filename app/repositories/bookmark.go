package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

type BookmarkRep struct {
	*mongoRep
}

func NewBookmarkRep(con *mongodb.Connection) *BookmarkRep {
	return &BookmarkRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Bookmark{}, con),
	}
}
