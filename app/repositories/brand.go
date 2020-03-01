package repositories

import (
	"go-shop-v2/pkg/repository"
)


type BrandRep struct {
	repository.IRepository
}

func NewBrandRep(rep repository.IRepository) *BrandRep {
	return &BrandRep{rep}
}
