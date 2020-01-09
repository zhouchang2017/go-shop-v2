package models

import (
	"go-shop-v2/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	register(func() *Inventory {
		return &Inventory{}
	})
}

type InventoryStatus int8

func (i InventoryStatus) Make(status int8) (res InventoryStatus) {
	switch status {
	case int8(ITEM_LOCKED):
		res = ITEM_LOCKED
	case int8(ITEM_OK):
		res = ITEM_OK
	case int8(ITEM_BAD):
		res = ITEM_BAD
	default:
		res = ITEM_PENDING
	}
	return res
}

const (
	ITEM_PENDING InventoryStatus = iota
	ITEM_LOCKED
	ITEM_OK
	ITEM_BAD
)

// 库存
type Inventory struct {
	model.MongoModel `inline`
	Shop             *AssociatedShop `json:"shop" bson:"shop"` // 门店
	Item             *AssociatedItem `json:"item"`             // sku
	Qty              int64           `json:"qty"`              // 存量
	Status           InventoryStatus `json:"status"`           // 状态
}

func (I Inventory) StatusOkMap() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":  "良品",
			"value": ITEM_OK,
		},
		{
			"name":  "不良品",
			"value": ITEM_BAD,
		},
	}
}

func (I Inventory) StatusMap() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":  "等待确认",
			"value": ITEM_PENDING,
		},
		{
			"name":  "锁定",
			"value": ITEM_LOCKED,
		},
		{
			"name":  "良品",
			"value": ITEM_OK,
		},
		{
			"name":  "不良品",
			"value": ITEM_BAD,
		},
	}
}

// 门店库存聚合结构
type AggregateShopCountStockInventory struct {
	ShopId   string                                    `json:"shop_id" bson:"shop_id"`
	ShopName string                                    `json:"shop_name" bson:"shop_name"`
	Total    int64                                     `json:"total"`
	Status   []*AggregateShopCountStockInventoryStatus `json:"status"`
}

type AggregateShopCountStockInventoryStatus struct {
	Status InventoryStatus `json:"status"`
	Qty    int64           `json:"qty"`
}

// 聚合结构体
type AggregateInventory struct {
	*AssociatedItem `inline`
	Total           int64                     `json:"total"` // 总数
	Inventories     []*AggregateInventoryUnit `json:"inventories"`
}

// 聚合库存单元
type AggregateInventoryUnit struct {
	Qty    int64                `json:"qty"`                // 存量
	Status InventoryStatus      `json:"status"`             // 状态
	Shops  []*InventoryUnitShop `json:"shops" bson:"shops"` // 门店
}

// 聚合门店结构
type InventoryUnitShop struct {
	Id          string             `json:"id"`
	InventoryId primitive.ObjectID `json:"inventory_id" bson:"inventory_id"`
	Name        string             `json:"name"`
	Qty         int64              `json:"qty"` // 存量
}

func (this *AggregateInventoryUnit) findShop(id string) (*InventoryUnitShop, bool) {
	for _, shop := range this.Shops {
		if shop.Id == id {
			return shop, true
		}
	}
	return nil, false
}

// 填充0值门店
func (this *AggregateInventory) WarpShops(shops []*InventoryUnitShop) *AggregateInventory {

	for _, unit := range this.Inventories {

		var unitShops []*InventoryUnitShop
		for _, shop := range shops {

			if unitShop, b := unit.findShop(shop.Id); b {
				unitShops = append(unitShops, unitShop)
				continue
			}
			unitShops = append(unitShops, shop)

		}
		unit.Shops = unitShops
	}

	return this
}

func (this *Inventory) SetStatus(status int8) {
	this.Status = this.Status.Make(status)
}
