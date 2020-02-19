package models

import "go-shop-v2/pkg/db/model"

// 购物车
type ShopCart struct {
	model.MongoModel `inline`
	UserId           string `json:"user_id" bson:"user_id"`
}

type ShopCartItem struct {
	ProductId    string         `json:"product_id" bson:"product_id"`                           // 产品id
	ItemId       string         `json:"item_id" bson:"item_id"`                                 // sku id
	OptionValues []*OptionValue `json:"option_values" bson:"option_values" form"option_values"` // 销售属性值
	Qty          int64          `json:"qty"`                                                    // 数量
	Checked      bool           `json:"checked"`                                                // 用户是否选定
}
