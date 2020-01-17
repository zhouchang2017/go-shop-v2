package models

import (
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db/model"
)

type InventoryChanger interface {
	OriginName() string
	OriginModel() string
	OriginId() string
}

type InventoryChange struct {
	Name  string `json:"name"`
	Model string `json:"model"`
	Id    string `json:"id"`
}

// 库存变动日志
type InventoryLog struct {
	model.MongoModel `inline`
	InventoryId      string           `json:"inventory_id" bson:"inventory_id"`
	ShopId           string           `json:"shop_id" bson:"shop_id"` // 冗余，用于计算门店每日出库量、入库量
	ItemId           string           `json:"item_id" bson:"item_id"` // 冗余，用于统计商品的库存变动
	BeforeQty        int64            `json:"before_qty" bson:"before_qty"`
	AfterQty         int64            `json:"after_qty" bson:"after_qty"`
	Value            int64            `json:"value"`  // 数量
	Origin           *InventoryChange `json:"origin"` // 来源
	User             *AssociatedAdmin `json:"user"`   // 操作用户
}

func NewInventoryHistory(inventory *Inventory, qty int64, origin InventoryChanger, user auth.Authenticatable) *InventoryLog {
	log := &InventoryLog{
		Origin: &InventoryChange{
			Name:  origin.OriginName(),
			Model: origin.OriginModel(),
			Id:    origin.OriginId(),
		},
		InventoryId: inventory.GetID(),
		ShopId:      inventory.Shop.Id,
		ItemId:      inventory.Item.Id,
		BeforeQty:   inventory.Qty - qty,
		AfterQty:    inventory.Qty,
		Value:       qty,
	}
	if admin, ok := user.(*Admin); ok {
		log.User = admin.ToAssociated()
	}
	return log
}
