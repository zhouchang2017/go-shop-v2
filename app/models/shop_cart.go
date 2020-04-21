package models

import (
	"go-shop-v2/pkg/db/model"
)

// 购物车
type ShopCart struct {
	model.MongoModel `inline`
	UserId           string          `json:"user_id" bson:"user_id"`
	Items            []*ShopCartItem `json:"items"`
}

type ShopCartItem struct {
	ItemId     string           `json:"item_id" bson:"item_id"`
	Item       *Item            `json:"item" bson:"-"`
	Promotions []*PromotionItem `json:"promotions" bson:"-"` // 促销信息
	Price      uint64            `json:"price"`               // 加入购物车时候的价格
	Qty        uint64            `json:"qty"`
}
