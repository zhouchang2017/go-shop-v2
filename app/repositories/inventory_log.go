package repositories

import (
	"go-shop-v2/pkg/repository"
)

type InventoryLogRep struct {
	repository.IRepository
}

func NewInventoryLogRep(rep repository.IRepository) *InventoryLogRep {
	return &InventoryLogRep{rep}
}
