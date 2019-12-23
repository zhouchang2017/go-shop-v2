package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

func init() {
	register(NewManualInventoryAction)
}

type ManualInventoryActionRep struct {
	*mongoRep
}

func NewManualInventoryAction(con *mongodb.Connection) *ManualInventoryActionRep {
	return &ManualInventoryActionRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.ManualInventoryAction{}, con),
	}
}
