package models

import "go-shop-v2/pkg/db/model"

// 购物车
type ShopCart struct {
	model.MongoModel `inline`
	UserId           string `json:"user_id" bson:"user_id"`
}

type ShopCartItem struct {
	Item    *AssociatedItem `json:"item" bson:"item"` // sku id
	Qty     int64           `json:"qty"`              // 数量
	Checked bool            `json:"checked"`          // 用户是否选定
}
