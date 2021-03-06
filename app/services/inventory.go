package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func init() {
	register(NewInventoryService)
}

type InventoryService struct {
	rep       *repositories.InventoryRep
	shopRep   *repositories.ShopRep
	itemRep   *repositories.ItemRep
	actionRep *repositories.ManualInventoryActionRep
}

func (this *InventoryService) GetRepository() *repositories.InventoryRep {
	return this.rep
}

func NewInventoryService(rep *repositories.InventoryRep, shopRep *repositories.ShopRep, itemRep *repositories.ItemRep, actionRep *repositories.ManualInventoryActionRep) *InventoryService {
	return &InventoryService{rep: rep, shopRep: shopRep, itemRep: itemRep, actionRep: actionRep}
}

func (this *InventoryService) Aggregate(ctx context.Context, req *request.IndexRequest) (data []*models.AggregateInventory, pagination response.Pagination, err error) {
	// 获取所有门店
	var shops []*models.InventoryUnitShop
	shopsRes := <-this.shopRep.FindAll(ctx)
	if shopsRes.Error != nil {
		return nil, pagination, shopsRes.Error
	}
	s := shopsRes.Result.([]*models.Shop)
	for _, shop := range s {
		shops = append(shops, &models.InventoryUnitShop{
			Id:   shop.GetID(),
			Name: shop.Name,
			Qty:  0,
		})
	}
	req.SetSearchField("item.code")
	aggregateRes := <-this.rep.AggregatePagination(ctx, req)
	if aggregateRes.Error != nil {
		return nil, aggregateRes.Pagination, aggregateRes.Error
	}
	data = aggregateRes.Result.([]*models.AggregateInventory)
	for _, item := range data {
		item.WarpShops(shops)
	}
	return data, aggregateRes.Pagination, nil
}

// 入库
func (this *InventoryService) Put(ctx context.Context, shopId string, itemId string, qty int64, status int8) (inventory *models.Inventory, err error) {
	// 检查当前是否存在对应规格产品库存
	incQtyRes := <-this.rep.IncQty(ctx, bson.M{
		"shop.id": shopId,
		"item.id": itemId,
		"status":  status,
	}, qty)
	if incQtyRes.Error == nil {
		return incQtyRes.Result.(*models.Inventory), nil
	}
	// 新增记录
	shopRes := <-this.shopRep.FindById(ctx, shopId)
	if shopRes.Error != nil {
		return nil, shopRes.Error
	}
	itemRes := <-this.itemRep.FindById(ctx, itemId)
	if itemRes.Error != nil {
		return nil, itemRes.Error
	}
	shop := shopRes.Result.(*models.Shop).ToAssociated()
	item := itemRes.Result.(*models.Item).ToAssociated()
	inventory = &models.Inventory{
		Shop: shop,
		Item: item,
		Qty:  qty,
	}
	inventory.SetStatus(status)
	createdRes := <-this.rep.Create(ctx, inventory)
	if createdRes.Error != nil {
		return nil, createdRes.Error
	}
	inventory = createdRes.Result.(*models.Inventory)
	return inventory, nil
}

// 出库
func (this *InventoryService) Take(ctx context.Context, id string, qty int64) (inventory *models.Inventory, err error) {
	byIdRes := <-this.rep.FindById(ctx, id)
	if byIdRes.Error != nil {
		return nil, byIdRes.Error
	}
	inventory = byIdRes.Result.(*models.Inventory)
	if inventory.Qty >= qty {
		incQtyRes := <-this.rep.IncQty(ctx, bson.M{"_id": inventory.ID}, -qty)
		if incQtyRes.Error != nil {
			return nil, incQtyRes.Error
		}
		inventory = incQtyRes.Result.(*models.Inventory)
		return inventory, nil
	}
	if bytes, err := json.Marshal(inventory); err == nil {
		log.Printf("剩余库存不足！剩余库存 %d ,需出库数量 %d，inventory:%s\n", inventory.Qty, qty, bytes)
	}

	return nil, fmt.Errorf("剩余库存不足！剩余库存 %d ,需出库数量 %d", inventory.Qty, qty)
}

// 手动操作记录
func (this *InventoryService) ManualActions(ctx context.Context, req *request.IndexRequest) <-chan repository.QueryPaginationResult {

	return this.actionRep.Pagination(ctx, req)
}
