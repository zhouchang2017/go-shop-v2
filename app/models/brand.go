package models

import "go-shop-v2/pkg/db/model"

func init() {
	register(NewBrand)
}

type Brand struct {
	model.MongoModel `inline`
	Name             string `json:"name" bson:"name"`
}

func NewBrand() *Brand {
	return &Brand{}
}


// 关联 简单brand
type AssociatedBrand struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
