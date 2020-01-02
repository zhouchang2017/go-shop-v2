package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
)

func init() {
	register(NewManualInventoryActionService)
}

type ManualInventoryActionService struct {
	rep     *repositories.ManualInventoryActionRep
	shopRep *repositories.ShopRep
	itemRep *repositories.ItemRep
}

func NewManualInventoryActionService(rep *repositories.ManualInventoryActionRep, shopRep *repositories.ShopRep, itemRep *repositories.ItemRep) *ManualInventoryActionService {
	return &ManualInventoryActionService{rep: rep, shopRep: shopRep, itemRep: itemRep}
}

type PutActionItem interface {
	GetId() string
	GetQty() int64
}

type InventoryActionPutOption struct {
	ShopId string                       `json:"shop_id" form:"shop_id"`
	Items  []*InventoryActionItemOption `json:"items" form:"items"`
}

func (this *ManualInventoryActionService) Put(ctx *gin.Context, option *InventoryActionPutOption, user *models.Admin) (*models.ManualInventoryAction, error) {
	action := &models.ManualInventoryAction{}
	action.SetTypeToPut()
	// 创建设置为保存状态
	action.SetStatusToSaved()
	action.User = user.ToAssociated()

	if _, err := this.SetShop(ctx, action, option.ShopId); err != nil {
		return action, err
	}
	if _, err := this.SetItems(ctx, action, option.Items...); err != nil {
		return action, err
	}
	created := <-this.rep.Create(ctx, action)
	if created.Error != nil {
		return action, created.Error
	}
	return created.Result.(*models.ManualInventoryAction), nil
}

type InventoryActionItemOption struct {
	Id          string `json:"id"`
	Qty         int64  `json:"qty"`
	InventoryId string `json:"inventory_id" form:"inventory_id"`
	Status      int8   `json:"status"`
}

type InventoryActionTakeOption struct {
	ShopId string                       `json:"shop_id" form:"shop_id"`
	Items  []*InventoryActionItemOption `json:"items"`
}

func (this *ManualInventoryActionService) Take(ctx *gin.Context, option *InventoryActionTakeOption, user *models.Admin)  (*models.ManualInventoryAction, error)  {
	action := &models.ManualInventoryAction{}
	action.SetTypeToTake()
	action.SetStatusToSaved()
	// TODO 检查库存，标记锁定状态
	if _, err := this.SetShop(ctx, action, option.ShopId); err != nil {
		return action, err
	}
	if _, err := this.SetItems(ctx, action, option.Items...); err != nil {
		return action, err
	}
	created := <-this.rep.Create(ctx, action)
	if created.Error != nil {
		return action, created.Error
	}
	return created.Result.(*models.ManualInventoryAction), nil
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

func (this *ManualInventoryActionService) SetItems(ctx context.Context, entity *models.ManualInventoryAction, items ...*InventoryActionItemOption) (*models.ManualInventoryAction, error) {
	var itemIds []string
	var manualInventoryActionItems []*models.ManualInventoryActionItem
	itemMap := map[string]int64{}
	for _, item := range items {
		itemIds = append(itemIds, item.Id)
		itemMap[item.Id] = item.Qty
	}

	if len(itemIds) > 0 {
		result := <-this.itemRep.FindByIds(ctx, itemIds...)
		if result.Error != nil {
			return entity, result.Error
		}

		for _, item := range result.Result.([]*models.Item) {
			associated := item.ToAssociated()
			qty := itemMap[item.GetID()]
			actionItem := &models.ManualInventoryActionItem{
				AssociatedItem: associated,
				Qty:            qty,
			}
			actionItem.SetStatusToPending()
			manualInventoryActionItems = append(manualInventoryActionItems, actionItem)
		}
	}

	entity.Items = manualInventoryActionItems
	return entity, nil
}
