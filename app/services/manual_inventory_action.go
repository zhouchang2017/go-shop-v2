package services

import (
	"context"
	"fmt"
	"github.com/novalagung/gubrak"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"log"
)

type ManualInventoryActionService struct {
	rep              *repositories.ManualInventoryActionRep
	inventoryService *InventoryService
	shopService      *ShopService
	productService   *ProductService
}

func NewManualInventoryActionService(
	rep *repositories.ManualInventoryActionRep,
	inventoryService *InventoryService,
	shopService *ShopService,
	productService *ProductService) *ManualInventoryActionService {
	return &ManualInventoryActionService{
		rep:              rep,
		inventoryService: inventoryService,
		shopService:      shopService,
		productService:   productService,
	}
}

type PutActionItem interface {
	GetId() string
	GetQty() int64
}

type InventoryActionPutOption struct {
	ShopId string                       `json:"shop_id" form:"shop_id"`
	Items  []*InventoryActionItemOption `json:"items" form:"items"`
}

// 入库单
func (this *ManualInventoryActionService) Put(ctx context.Context, option *InventoryActionPutOption, user *models.Admin) (*models.ManualInventoryAction, error) {
	action := &models.ManualInventoryAction{}
	action.SetTypeToPut()
	// 创建设置为保存状态
	action.SetStatusToSaved()
	action.User = user.ToAssociated()

	if _, err := this.setShop(ctx, action, option.ShopId); err != nil {
		return action, err
	}
	if _, err := this.setItems(ctx, action, option.Items...); err != nil {
		return action, err
	}
	created := <-this.rep.Create(ctx, action)
	if created.Error != nil {
		return action, created.Error
	}
	return created.Result.(*models.ManualInventoryAction), nil
}

// 更新入库单
func (this *ManualInventoryActionService) UpdatePut(ctx context.Context, id string, option *InventoryActionPutOption, user *models.Admin) (*models.ManualInventoryAction, error) {
	// 获取入库单
	actionResult := <-this.rep.FindById(ctx, id)
	if actionResult.Error != nil {
		return nil, actionResult.Error
	}
	action := actionResult.Result.(*models.ManualInventoryAction)
	// 创建设置为保存状态
	action.SetStatusToSaved()
	// 设置更新用户
	action.User = user.ToAssociated()
	// 设置门店
	if _, err := this.setShop(ctx, action, option.ShopId); err != nil {
		return action, err
	}
	// 设置商品集
	if _, err := this.setItems(ctx, action, option.Items...); err != nil {
		return action, err
	}
	saved := <-this.rep.Save(ctx, action)
	if saved.Error != nil {
		return action, saved.Error
	}
	return saved.Result.(*models.ManualInventoryAction), nil
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

// 出库单
func (this *ManualInventoryActionService) Take(ctx context.Context, option *InventoryActionTakeOption, user *models.Admin) (*models.ManualInventoryAction, error) {
	action := &models.ManualInventoryAction{}
	action.SetTypeToTake()
	action.SetStatusToSaved()
	action.User = user.ToAssociated()

	if _, err := this.setShop(ctx, action, option.ShopId); err != nil {
		return action, err
	}
	if _, err := this.setItems(ctx, action, option.Items...); err != nil {
		return action, err
	}
	// 标记锁定状态,TODO 开启事务
	for _, item := range action.Items {
		if err := this.inventoryService.LockById(ctx, item.InventoryId, item.Qty); err != nil {
			return action, err
		}
	}

	created := <-this.rep.Create(ctx, action)
	if created.Error != nil {
		return action, created.Error
	}

	return created.Result.(*models.ManualInventoryAction), nil
}

// 更新出库单
func (this *ManualInventoryActionService) UpdateTake(ctx context.Context, id string, option *InventoryActionTakeOption, user *models.Admin) (*models.ManualInventoryAction, error) {
	// 获取出库单
	actionResult := <-this.rep.FindById(ctx, id)
	if actionResult.Error != nil {
		return nil, actionResult.Error
	}
	action := actionResult.Result.(*models.ManualInventoryAction)
	// 创建设置为保存状态
	action.SetStatusToSaved()
	// 设置更新用户
	action.User = user.ToAssociated()

	if action.Shop.Id != option.ShopId {
		if _, err := this.setShop(ctx, action, option.ShopId); err != nil {
			return action, err
		}
	}

	// TODO 开启事务
	// 解锁之前锁定库存
	for _, item := range action.Items {
		err := this.inventoryService.UnLockById(ctx, item.InventoryId, item.Qty)
		if err != nil {
			return action, err
		}
	}
	// 设置新的items
	if _, err := this.setItems(ctx, action, option.Items...); err != nil {
		return action, err
	}
	// 标记锁定状态
	for _, item := range action.Items {
		if err := this.inventoryService.LockById(ctx, item.InventoryId, item.Qty); err != nil {
			return action, err
		}
	}

	saved := <-this.rep.Save(ctx, action)

	if saved.Error != nil {
		return action, saved.Error
	}
	return saved.Result.(*models.ManualInventoryAction), nil
}

// 确认操作
func (this *ManualInventoryActionService) StatusToFinished(ctx context.Context, id string) (entity *models.ManualInventoryAction, err error) {
	action, err := this.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if action.Type.IsPut() {
		// 入库
		if err := this.putInventory(ctx, action); err != nil {
			return nil, err
		}
	} else {
		// 出库
		if err := this.takeInventory(ctx, action); err != nil {
			return nil, err
		}
	}

	action.SetStatusToFinished()

	saved := <-this.rep.Save(ctx, action)

	if saved.Error != nil {
		return nil, saved.Error
	}
	return saved.Result.(*models.ManualInventoryAction), nil
}

// 取消操作
func (this *ManualInventoryActionService) Cancel(ctx context.Context, id string) (err error) {
	action, err := this.FindById(ctx, id)
	if err != nil {
		return err
	}

	if action.Type.IsTake() {
		// 出库单，归还锁定库存
		for _, item := range action.Items {
			if err := this.inventoryService.UnLock(ctx, action.Shop.Id, item.Id, item.Qty, int8(item.Status)); err != nil {
				return err
			}
		}
	}

	action.SetStatusToCancel()

	saved := <-this.rep.Save(ctx, action)

	return saved.Error
}

// 推入仓库
func (this *ManualInventoryActionService) putInventory(ctx context.Context, action *models.ManualInventoryAction) error {

	for _, item := range action.Items {
		_, err := this.inventoryService.Put(ctx, action.Shop.Id, item.Id, uint64(item.Qty), int8(item.Status), action)
		if err != nil {
			return err
		}
	}

	return nil
}

// 推出仓库
func (this *ManualInventoryActionService) takeInventory(ctx context.Context, action *models.ManualInventoryAction) error {
	// TODO 事物
	for _, item := range action.Items {
		if err := this.inventoryService.TakeByLocked(ctx, item.InventoryId, item.Qty, action); err != nil {
			return err
		}
	}
	return nil
}

func (this *ManualInventoryActionService) setShop(ctx context.Context, entity *models.ManualInventoryAction, shopId string) (*models.ManualInventoryAction, error) {
	shop, err := this.shopService.FindById(ctx, shopId)
	if err != nil {
		return nil, err
	}
	entity.Shop = shop.ToAssociated()
	return entity, nil
}

func (this *ManualInventoryActionService) setItems(ctx context.Context, entity *models.ManualInventoryAction, items ...*InventoryActionItemOption) (*models.ManualInventoryAction, error) {
	var itemIds []string
	var manualInventoryActionItems []*models.ManualInventoryActionItem
	for _, item := range items {
		itemIds = append(itemIds, item.Id)
	}

	itemIds = gubrak.From(itemIds).Uniq().Result().([]string)

	var productItems []*models.Item

	if len(itemIds) > 0 {
		// chunk
		gubrak.From(itemIds).Chunk(50).Each(func(value []string) {
			result := <-this.productService.itemRep.FindByIds(ctx, value)
			if result.Error != nil {
				log.Printf("findByIds error:%s\n", result.Error)
			}

			productItems = append(productItems, result.Result.([]*models.Item)...)
		})

		for _, item := range items {
			var currentItem *models.Item

			currentItem = gubrak.From(productItems).Find(func(i *models.Item) bool {
				return i.GetID() == item.Id
			}).Result().(*models.Item)
			actionItem := &models.ManualInventoryActionItem{
				AssociatedItem: currentItem.ToAssociated(),
				Qty:            item.Qty,
			}
			actionItem.SetStatus(item.Status)
			if item.InventoryId != "" {
				actionItem.InventoryId = item.InventoryId
			}

			manualInventoryActionItems = append(manualInventoryActionItems, actionItem)
		}
	}

	entity.Items = manualInventoryActionItems
	return entity, nil
}

// 操作详情附加库存数据
func (this *ManualInventoryActionService) FindByIdWithInventory(ctx context.Context, id string) (*models.ManualInventoryAction, error) {
	result := <-this.rep.FindById(ctx, id)
	if result.Error != nil {
		return nil, result.Error
	}

	action := result.Result.(*models.ManualInventoryAction)

	var inventoryIds []string
	for _, item := range action.Items {
		if item.InventoryId != "" {
			inventoryIds = append(inventoryIds, item.InventoryId)
		}
	}
	if len(inventoryIds) > 0 {
		inventories, err := this.inventoryService.FindByIds(ctx, inventoryIds...)

		if err != nil {
			return nil, err
		}

		for _, item := range action.Items {
			inventory, err := this.resolveInventoryById(inventories, item.InventoryId)
			if err != nil {
				return nil, err
			}
			item.Inventory = inventory
		}
	}
	return action, nil
}

// 列表页
func (this *ManualInventoryActionService) Pagination(ctx context.Context, req *request.IndexRequest) (action []*models.ManualInventoryAction, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.ManualInventoryAction), results.Pagination, nil
}

// 详情页
func (this *ManualInventoryActionService) FindById(ctx context.Context, id string) (action *models.ManualInventoryAction, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}
	return byId.Result.(*models.ManualInventoryAction), nil
}

func (this *ManualInventoryActionService) FindByIds(ctx context.Context, ids ...string) (action []*models.ManualInventoryAction, err error) {
	results := <-this.rep.FindByIds(ctx, ids)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.ManualInventoryAction), nil
}

func (this *ManualInventoryActionService) resolveInventoryById(inventories []*models.Inventory, id string) (*models.Inventory, error) {
	for _, inventory := range inventories {
		if inventory.GetID() == id {
			return inventory, nil
		}
	}
	return nil, fmt.Errorf("inventory id = %s,not found!!", id)
}
