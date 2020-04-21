package repositories

import (
	"go-shop-v2/pkg/repository"
)

type CategoryRep struct {
	repository.IRepository
}


func NewCategoryRep(rep repository.IRepository) *CategoryRep {
	return &CategoryRep{rep}
}
