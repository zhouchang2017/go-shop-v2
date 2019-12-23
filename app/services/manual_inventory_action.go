package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
)

func init() {
	register(NewManualInventoryActionService)
}

type ManualInventoryActionService struct {
	rep     *repositories.ManualInventoryActionRep
	shopRep *repositories.ShopRep
}

func NewManualInventoryActionService(rep *repositories.ManualInventoryActionRep, shopRep *repositories.ShopRep) *ManualInventoryActionService {
	return &ManualInventoryActionService{rep: rep, shopRep: shopRep}
}

func (this *ManualInventoryActionService) SetShop(ctx context.Context, entity *models.ManualInventoryAction, shopId string) (*models.ManualInventoryAction, error) {
	shopRes := <-this.shopRep.FindById(ctx, shopId)
	if shopRes.Error != nil {
		return nil, shopRes.Error
	}
	shop := shopRes.Result.(*models.Shop)
	entity.Shop = shop.ToAssociated()
	return entity, nil
}
