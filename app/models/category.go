package models

import (
	uuid "github.com/satori/go.uuid"
	"go-shop-v2/pkg/db/model"
)


type Category struct {
	model.MongoModel `inline`
	Name             string           `json:"name"`
	Options          []*ProductOption `json:"options"`
}

func NewCategory(name string) *Category {
	return &Category{Name: name, Options: []*ProductOption{}}
}

type AssociatedCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// 产品销售属性
type ProductOption struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	Sort   int64          `json:"sort"`
	Values []*OptionValue `json:"values"`
}

func (this *ProductOption) NewValue(value, code string) *OptionValue {
	return &OptionValue{
		Code:  code,
		Value: value,
	}
}

func (this *ProductOption) AddValues(values ...*OptionValue) *ProductOption {
	for _, value := range values {
		if optionValue, b := this.findValueByCode(value.Code); b {
			optionValue.Value = value.Value
		} else {
			this.Values = append(this.Values, value)
		}
	}
	return this
}

func (this *ProductOption) findValueByCode(code string) (*OptionValue, bool) {
	for _, v := range this.Values {
		if v.Code == code {
			return v, true
		}
	}
	return nil, false
}

// 产品销售属性值
type OptionValue struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}


func (this *OptionValue) Equal(value *OptionValue) bool {
	return this.Code == value.Code
}

func NewProductOption(name string) *ProductOption {
	return &ProductOption{
		Id:     uuid.NewV4().String(),
		Name:   name,
		Values: []*OptionValue{},
	}
}

func MakeProductOption(id string, name string, sort int64) *ProductOption {
	return &ProductOption{
		Id:   id,
		Name: name,
		Sort: sort,
	}
}
