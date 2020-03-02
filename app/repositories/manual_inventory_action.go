package repositories

import (
	"go-shop-v2/pkg/repository"
)


type ManualInventoryActionRep struct {
	repository.IRepository
}

func NewManualInventoryActionRep(rep repository.IRepository) *ManualInventoryActionRep {
	return &ManualInventoryActionRep{rep}
}
