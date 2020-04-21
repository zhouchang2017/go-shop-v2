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
	Attributes       []*ProductAttribute `json:"attributes,omitempty" bson:"attributes"`
	Options          []*ProductOption    `json:"options,omitempty" bson:"options"`
	Items            []*Item             `json:"items,omitempty" bson:"-"`
	Description      string              `json:"description,omitempty" bson:"description"`
	Price            uint64               `json:"price" bson:"price"`
	PromotionPrice   uint64               `json:"promotion_price" bson:"-"` // 促销价，显示items 最低价格
	TotalSalesQty    uint64               `json:"total_sales_qty" bson:"total_sales_qty"`
	FakeSalesQty     uint64               `json:"fake_sales_qty,omitempty" bson:"fake_sales_qty"`
	Images           []qiniu.Image       `json:"images,omitempty" bson:"images"`
	Avatar           *qiniu.Image        `json:"avatar,omitempty" bson:"avatar"`
	OnSale           bool                `json:"on_sale" bson:"on_sale"`
	Sort             int64               `json:"sort"` // 排序权重
	Qty              uint64               `json:"qty" bson:"-"`
	CollectCount     uint64               `json:"collect_count" bson:"collect_count"` // 累计收藏
	ShareCount       uint64               `json:"share_count" bson:"share_count"`     // 累计分享
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
	Price    uint64               `json:"price"`
}

func (this Product) GetSort() int64 {
	return this.Sort
}

func (this Product) GetType() string {
	return "product"
}

func (this *Product) SetAvatar() {
	var avatar *qiniu.Image
	if len(this.Images) > 0 {
		avatar = &(this.Images[0])
	}
	this.Avatar = avatar
}

func (this Product) GetAvatar() *qiniu.Image {
	return this.Avatar
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
		Price:    this.Price,
	}
}

// 产品属性
type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
