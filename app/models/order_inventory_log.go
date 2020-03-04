package models

import "go-shop-v2/pkg/db/model"

type OrderInventoryLog struct {
	model.MongoModel
	OrderNo          string              `json:"order_no" bson:"order_no" name:"订单号"`
	ItemInventoryLog []*ItemInventoryLog `json:"item_inventory_log" bson:"item_inventory_log" name:"商品库存发货记录"`
}

type ItemInventoryLog struct {
	ItemId      string `json:"item_id" bson:"item_id"`
	InventoryId string `json:"inventory_id" bson:"inventory_id"`
}
