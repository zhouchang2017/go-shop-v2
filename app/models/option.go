package models

import (
	uuid "github.com/satori/go.uuid"
	"go-shop-v2/pkg/qiniu"
)

// 产品销售属性
type ProductOption struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	Image  bool           `json:"image"`
	Values []*OptionValue `json:"values"`
}

type ProductOptions []*ProductOption

func (p ProductOptions) Len() int {
	return len(p)
}

func (p ProductOptions) Less(i, j int) bool {
	var iValue = 0
	var jValue = 0
	if p[i].Image {
		iValue = 1
	}
	if p[j].Image {
		jValue = 1
	}
	return iValue > jValue
}

func (p ProductOptions) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ProductOptions) Append(option *ProductOption) {
	p = append(p, option)
}

func (this *ProductOption) NewValue(name string) *OptionValue {
	return &OptionValue{
		Id:   uuid.NewV4().String(),
		Name: name,
	}
}

func (this *ProductOption) MakeValue(id string, name string) *OptionValue {
	return &OptionValue{
		Id:   id,
		Name: name,
	}
}

func (this *ProductOption) AddValues(values ...*OptionValue) *ProductOption {
	for _, value := range values {
		if optionValue, b := this.findValueById(value.Id); b {
			optionValue.Name = value.Name
			optionValue.Image = value.Image
		} else {
			this.Values = append(this.Values, value)
		}
	}
	return this
}

func (this *ProductOption) findValueById(id string) (*OptionValue, bool) {
	for _, v := range this.Values {
		if v.Id == id {
			return v, true
		}
	}
	return nil, false
}

// 产品销售属性值
type OptionValue struct {
	Id    string       `json:"id"`
	Name  string       `json:"name"`
	Image *qiniu.Image `json:"image,omitempty" bson:"image,omitempty"` // 缩略图
}

// 用于排序
type SortOptionValues struct {
	Values  []*OptionValue
	Options ProductOptions
}

func (o SortOptionValues) Len() int {
	return len(o.Values)
}

func (o SortOptionValues) Less(i, j int) bool {
	return o.getSort(i) > o.getSort(j)
}

func (o SortOptionValues) getSort(i int) int {
	id := o.Values[i].Id
	// 从options从获取排序值
	for index, opt := range o.Options {
		for _, value := range opt.Values {
			if value.Id == id {
				return len(o.Options) - index
			}
		}
	}
	return 0
}

func (o SortOptionValues) Swap(i, j int) {
	o.Values[i], o.Values[j] = o.Values[j], o.Values[i]
}

func (this *OptionValue) Equal(value *OptionValue) bool {
	return this.Id == value.Id
}

func (this *OptionValue) SetImage(url string) *OptionValue {
	if url == "" {
		this.Image = nil
		return this
	}
	image := qiniu.NewImage(url)
	this.Image = &image
	return this
}

func NewProductOption(name string) *ProductOption {
	return &ProductOption{
		Id:     uuid.NewV4().String(),
		Name:   name,
		Values: []*OptionValue{},
	}
}

func MakeProductOption(id string, name string) *ProductOption {
	return &ProductOption{
		Id:   id,
		Name: name,
	}
}
