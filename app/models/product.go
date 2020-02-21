package models

import (
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/qiniu"
)

// 产品
type Product struct {
	model.MongoModel `inline`
	Name             string              `json:"name" bson:"name"`
	Code             string              `json:"code" bson:"code"`
	Brand            *AssociatedBrand    `json:"brand" bson:"brand"`
	Category         *AssociatedCategory `json:"category" bson:"category"`
	Attributes       []*ProductAttribute `json:"attributes" bson:"attributes"`
	Options          []*ProductOption    `json:"options" bson:"options"`
	Items            []*Item             `json:"items,omitempty" bson:"-"`
	Description      string              `json:"description" bson:"description"`
	Price            int64               `json:"price" bson:"price"`
	TotalSalesQty    int64               `json:"total_sales_qty" bson:"total_sales_qty"`
	FakeSalesQty     int64               `json:"fake_sales_qty" bson:"fake_sales_qty"`
	Images           []qiniu.Image       `json:"images" bson:"images"`
	OnSale           bool                `json:"on_sale" bson:"on_sale"`
	Sort             int64               `json:"sort"`
	Qty              int64               `json:"qty" bson:"-"`
}

func NewProduct() *Product {
	return &Product{}
}

// 关联简单产品结构
type AssociatedProduct struct {
	Id       string              `json:"id"`
	Name     string              `json:"name"`
	Code     string              `json:"code"`
	Brand    *AssociatedBrand    `json:"brand"`
	Category *AssociatedCategory `json:"category"`
	Avatar   *qiniu.Image        `json:"avatar"`
}

func (this Product) GetSort() int64 {
	return this.Sort
}

func (this Product) GetType() string {
	return "product"
}

func (this Product) ToAssociated() *AssociatedProduct {
	var avatar *qiniu.Image
	if len(this.Images) > 0 {
		avatar = &(this.Images[0])
	}
	return &AssociatedProduct{
		Id:       this.GetID(),
		Name:     this.Name,
		Code:     this.Code,
		Brand:    this.Brand,
		Category: this.Category,
		Avatar:   avatar,
	}
}

// 产品属性
type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
