package models

import "go-shop-v2/pkg/db/model"

// 收藏夹
type Bookmark struct {
	model.MongoModel `inline`
	UserId           string   `json:"user_id" bson:"user_id"`
	ProductIds       []string `json:"product_ids" bson:"product_ids"`
}
