package models

import (
	"go-shop-v2/pkg/db/model"
)

type Category struct {
	model.MongoModel `inline`
	Name             string `json:"name"`
}

func NewCategory(name string) *Category {
	return &Category{Name: name}
}

type AssociatedCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
