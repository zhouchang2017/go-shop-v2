package repositories

import "go-shop-v2/pkg/repository"

type OrderInventoryLogRep struct {
	repository.IRepository
}

func NewOrderInventoryLogRep(rep repository.IRepository) *OrderInventoryLogRep {
	return &OrderInventoryLogRep{rep}
}
