package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

type InventoryService struct {
	rep            *repositories.InventoryRep
	historyRep     *repositories.InventoryLogRep
	shopService    *ShopService
	productService *ProductService
}

func (this *InventoryService) GetRepository() *repositories.InventoryRep {
	return this.rep
}

func NewInventoryService(rep *repositories.InventoryRep, historyRep *repositories.InventoryLogRep, shopService *ShopService, productService *ProductService) *InventoryService {
	return &InventoryService{rep: rep, historyRep: historyRep, shopService: shopService, productService: productService}
}

func (this *InventoryService) FindByIds(ctx context.Context, ids ...string) (inventories []*models.Inventory, err error) {
	byIds := <-this.rep.FindByIds(ctx, ids...)
	if byIds.Error != nil {
		err = byIds.Error
		return
	}

	return byIds.Result.([]*models.Inventory), nil
}

func (this *InventoryService) Aggregate(ctx context.Context, req *request.IndexRequest) (data []*models.AggregateInventory, pagination response.Pagination, err error) {
	// 获取所有门店
	var shops []*models.InventoryUnitShop
	shops2, err := this.shopService.All(ctx)
	if err != nil {
		return nil, pagination, err
	}
	for _, shop := range shops2 {
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

// 库存统计聚合
func (this *InventoryService) AggregateStockByShops(ctx context.Context, shopIds ...string) (data []*models.AggregateShopCountStockInventory, err error) {
	return this.rep.AggregateStockByShops(ctx, shopIds...)
}

// 入库
func (this *InventoryService) Put(ctx context.Context, shopId string, itemId string, qty int64, status int8, changer models.InventoryChanger) (inventory *models.Inventory, err error) {
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
	shop2, err := this.shopService.FindById(ctx, shopId)
	if err != nil {
		return nil, err
	}
	item2, err := this.productService.FindItemById(ctx, itemId)
	if err != nil {
		return nil, err
	}
	shop := shop2.ToAssociated()
	item := item2.ToAssociated()
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

	// 记录操作日志
	go this.writeLog(ctx, inventory, qty, changer)

	return inventory, nil
}

// 写操作日志
func (this *InventoryService) writeLog(ctx context.Context, inventory *models.Inventory, qty int64, origin models.InventoryChanger) {
	var user auth.Authenticatable
	authUser := ctx2.GetUser(ctx)
	if authUserd, ok := authUser.(auth.Authenticatable); ok {
		user = authUserd
	}
	created := <-this.historyRep.Create(ctx, models.NewInventoryHistory(inventory, qty, origin, user))
	if created.Error != nil {
		log.Printf("write log error:%s\n", created.Error)
	}
}

// 出库
func (this *InventoryService) Take(ctx context.Context, id string, qty int64, changer models.InventoryChanger) (inventory *models.Inventory, err error) {
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

	// 记录操作日志
	go this.writeLog(ctx, inventory, -qty, changer)

	return nil, fmt.Errorf("剩余库存不足！剩余库存 %d ,需出库数量 %d", inventory.Qty, qty)
}

// 锁定库存
func (this *InventoryService) Lock(ctx context.Context, shopId string, itemId string, qty int64, status int8) (err error) {
	return this.rep.Lock(ctx, shopId, itemId, qty, status)
}

// 通过id锁定库存
func (this *InventoryService) LockById(ctx context.Context, id string, qty int64) (err error) {
	return this.rep.LockById(ctx, id, qty)
}

// 解锁
func (this *InventoryService) UnLock(ctx context.Context, shopId string, itemId string, qty int64, status int8) (err error) {
	return this.rep.UnLock(ctx, shopId, itemId, qty, status)
}

// 通过id解锁
func (this *InventoryService) UnLockById(ctx context.Context, id string, qty int64) (err error) {
	return this.rep.UnLockById(ctx, id, qty)
}

// 通过锁定库存出库
func (this *InventoryService) TakeByLocked(ctx context.Context, id string, qty int64, changer models.InventoryChanger) (err error) {
	err = this.rep.TakeByLocked(ctx, id, qty)
	if err != nil {
		return err
	}
	// 记录操作日志
	go func() {
		inventory, err := this.FindById(ctx, id)
		if err != nil {
			log.Printf("find by id error:%s\n", err)
			return
		}
		this.writeLog(ctx, inventory, -qty, changer)
	}()

	return nil
}

// 列表
func (this *InventoryService) Pagination(ctx context.Context, req *request.IndexRequest) (inventories []*models.Inventory, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Inventory), results.Pagination, nil
}

// 详情
func (this *InventoryService) FindById(ctx context.Context, id string) (inventory *models.Inventory, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}
	return byId.Result.(*models.Inventory), nil
}

// 库存查询
func (this *InventoryService) Search(ctx context.Context, opt *repositories.QueryOption) (inventories []*models.Inventory, err error) {
	search := <-this.rep.Search(ctx, opt)
	return search.Result, search.Error
}

// 通过sku id获取库存
func (this *InventoryService) SearchItemsQty(ctx context.Context, ids ...string) (results map[string]int64) {
	var g errgroup.Group
	var mu sync.Mutex
	results = map[string]int64{}
	sem := make(chan struct{}, 10)
	for _, id := range ids {
		id := id
		sem <- struct{}{}
		g.Go(func() error {
			count, _ := this.rep.SearchByItemId(ctx, id)
			mu.Lock()
			results[id] = count
			mu.Unlock()
			<-sem
			return nil
		})
	}
	g.Wait()
	return results
}
