package models

import (
	"bytes"
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/qiniu"
)

type Item struct {
	model.MongoModel `inline`
	Code             string             `json:"code" bson:"code"`
	Product          *AssociatedProduct `json:"product,omitempty" bson:"product,omitempty"`
	Price            int64              `json:"price,omitempty" bson:"price"`
	PromotionPrice   int64              `json:"promotion_price" bson:"-"` // 促销价
	OptionValues     []*OptionValue     `json:"option_values" bson:"option_values" form"option_values" `
	SalesQty         int64              `json:"sales_qty,omitempty" bson:"sales_qty" form"sales_qty" `
	Qty              int64              `json:"qty" bson:"qty"` // 可售数量
	Inventories      []*Inventory       `json:"inventories,omitempty" bson:"-"`
	Avatar           *qiniu.Image       `json:"avatar,omitempty" bson:"avatar"`
}

func (this *Item) SetAvatar() {
	var image *qiniu.Image
	if this.Product != nil {
		image = this.Product.Avatar
	}
	for _, value := range this.OptionValues {
		if value.Image != nil && value.Image.Src() != "" {
			image = value.Image
		}
	}
	this.Avatar = image
}

func (this Item) GetAvatar() (image *qiniu.Image) {
	return this.Avatar
}

// 关联简单SKU结构
type AssociatedItem struct {
	Id           string             `json:"id"`
	Code         string             `json:"code"`                               // sku码
	Avatar       *qiniu.Image       `json:"avatar"`                             // 图
	Product      *AssociatedProduct `json:"product"`                            // 冗余产品信息
	Price        int64              `json:"price"`                              // 价格
	OptionValues []*OptionValue     `json:"option_values" bson:"option_values"` // sku销售属性
}

func NewItem() *Item {
	return &Item{}
}

func (this Item) ToAssociated() *AssociatedItem {
	return &AssociatedItem{
		Id:           this.GetID(),
		Code:         this.Code,
		Product:      this.Product,
		OptionValues: this.OptionValues,
		Avatar:       this.GetAvatar(),
		Price:        this.Price,
	}
}

func (this *Item) OptionValueString() string {
	bufferString := bytes.NewBufferString("")
	for _, opt := range this.OptionValues {
		bufferString.WriteString(opt.Name)
	}
	return bufferString.String()
}

// 添加销售属性值
func (this *Item) AddOptionValues(ov ...*OptionValue) {
	for _, value := range ov {
		if exist, _ := this.optionValueExist(value); exist {
			continue
		}
		this.OptionValues = append(this.OptionValues, value)
	}
}

// 销售属性值是否存在
func (this *Item) optionValueExist(ov *OptionValue) (exist bool, index int) {
	for index, value := range this.OptionValues {
		if value.Equal(ov) {
			return true, index
		}
	}
	return false, -1
}
