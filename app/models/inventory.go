package models

import "go-shop-v2/pkg/db/model"

func init() {
	register(func() *Inventory {
		return &Inventory{}
	})
}

type InventoryStatus int8

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


func (this *Inventory) SetStatus(status int8) {
	this.Status = InventoryStatus(status)
}
