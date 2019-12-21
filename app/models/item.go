package models

import (
	"bytes"
	"go-shop-v2/pkg/db/model"
)

func init() {
	register(NewItem)
}

type Item struct {
	model.MongoModel `inline`
	Code             string         `json:"code" bson:"code"`
	Product          *Product       `json:"product,omitempty" bson:"product,omitempty"`
	Price            int64          `json:"price,omitempty" bson:"price"`
	OptionValues     []*OptionValue `json:"option_values" bson:"option_values" form"option_values" `
	SalesQty         int64          `json:"sales_qty,omitempty" bson:"sales_qty" form"sales_qty" `
}

// 关联简单SKU结构
type AssociatedItem struct {
	Id           string             `json:"id"`
	Code         string             `json:"code"`                               // sku码
	Product      *AssociatedProduct `json:"product"`                            // 冗余产品信息
	OptionValues []*OptionValue     `json:"option_values" bson:"option_values"` // sku销售属性
}

func NewItem() *Item {
	return &Item{}
}

func (this Item) ToAssociated() *AssociatedItem {
	return &AssociatedItem{
		Id:           this.GetID(),
		Code:         this.Code,
		Product:      this.Product.ToAssociated(),
		OptionValues: this.OptionValues,
	}
}

func (this *Item) OptionValueString() string {
	bufferString := bytes.NewBufferString("")
	for _, opt := range this.OptionValues {
		bufferString.WriteString(opt.Value)
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
