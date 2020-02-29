package models

import "go-shop-v2/pkg/db/model"

// 收藏夹
type Bookmark struct {
	model.MongoModel `inline`
	UserId           string             `json:"user_id" bson:"user_id"`
	Product          *AssociatedProduct `json:"product" bson:"product"`
	Enabled          bool               `json:"enabled"` // 是否已失效,product 被删除后标记为失效
}
