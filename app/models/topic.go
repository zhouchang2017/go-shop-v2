package models

import (
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/qiniu"
)

// 专题
type Topic struct {
	model.MongoModel `inline`
	Title            string      `json:"title"`
	ShortTitle       string      `json:"short_title" bson:"short_title"`
	Avatar           qiniu.Image `json:"avatar"`
	Content          string      `json:"content"`
	ProductIds       []string    `json:"product_ids"`
	Sort             int64       `json:"sort"`
}

func (this Topic) GetSort() int64 {
	return this.Sort
}

func (this Topic) GetType() string {
	return "topic"
}
