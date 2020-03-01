package repositories

import (
	"go-shop-v2/pkg/repository"
)

type OrderRep struct {
	repository.IRepository
}

func NewOrderRep(rep repository.IRepository) *OrderRep {
	return &OrderRep{rep}
}

