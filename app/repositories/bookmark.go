package repositories

import (
	"go-shop-v2/pkg/repository"
)

type BookmarkRep struct {
	repository.IRepository
}

func NewBookmarkRep(rep repository.IRepository) *BookmarkRep {
	return &BookmarkRep{rep}
}
