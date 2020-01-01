package services

import (
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
)

func MakeInventoryService() *InventoryService {
	con := mongodb.GetConFn()
	rep := repositories.NewInventoryRep(con)
	shopRep := repositories.NewShopRep(con)
	itemRep := repositories.NewItemRep(con)
	actionRep := repositories.NewManualInventoryActionRep(con)
	return NewInventoryService(rep, shopRep, itemRep, actionRep)
}
