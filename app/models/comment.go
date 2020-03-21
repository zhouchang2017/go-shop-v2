package models

import (
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/qiniu"
)

// 商品评论
type Comment struct {
	model.MongoModel `inline`
	ProductId        string         `json:"product_id" bson:"product_id"`
	ItemId           string         `json:"item_id" bson:"item_id"`
	OrderId          string         `json:"order_id" bson:"order_id"`
	User             *CommentUser   `json:"user" bson:"user"`
	Content          string         `json:"content"` // 内容
	Images           []*qiniu.Image `json:"images"`
	Rate             float64        `json:"rate"` // 打分
}

type CommentUser struct {
	UserId   string `json:"user_id" bson:"user_id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}
