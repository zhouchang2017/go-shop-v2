package models

import (
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/qiniu"
)

// 文章
type Article struct {
	model.MongoModel `inline`
	Title            string        `json:"title"`
	ShortTitle       string        `json:"short_title" bson:"short_title"`
	Photos           []qiniu.Image `json:"photos"`
	Content          string        `json:"content,omitempty"`
	ProductId        string        `json:"product_id,omitempty" bson:"product_id"`
	Sort             int64         `json:"sort"`
}

func (this Article) GetSort() int64 {
	return this.Sort
}

func (this Article) GetType() string {
	return "article"
}
