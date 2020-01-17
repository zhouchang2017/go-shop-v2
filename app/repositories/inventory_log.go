package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

type InventoryLogRep struct {
	*mongoRep
}

func NewInventoryLogRep(con *mongodb.Connection) *InventoryLogRep {
	return &InventoryLogRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.InventoryLog{}, con),
	}
}
