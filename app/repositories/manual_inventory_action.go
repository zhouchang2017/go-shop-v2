package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)


type ManualInventoryActionRep struct {
	*mongoRep
}

func NewManualInventoryActionRep(con *mongodb.Connection) *ManualInventoryActionRep {
	return &ManualInventoryActionRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.ManualInventoryAction{}, con),
	}
}
