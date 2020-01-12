package models

import (
	"fmt"
	"go-shop-v2/pkg/db/model"
)

// 门店
type Shop struct {
	model.MongoModel `inline`
	Name             string             `json:"name" bson:"name"`
	Address          *ShopAddress       `json:"address"`
	Location         *Location          `json:"location"` // 坐标
	Members          []*AssociatedAdmin `json:"members"`  // 成员
}

// 关联简单门店结构
type AssociatedShop struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}

func (s Shop) ToAssociated() *AssociatedShop {

	res := &AssociatedShop{
		Id:   s.GetID(),
		Name: s.Name,
	}

	if s.Location != nil {
		res.Location = s.Location.GeoJSON()
	}

	return res
}

func NewShop() *Shop {
	return &Shop{}
}

// 门店地址
type ShopAddress struct {
	Addr     string `json:"addr"`
	Areas    string `json:"areas"`
	City     string `json:"city"`
	Province string `json:"province"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

func (this ShopAddress) String() string {
	return fmt.Sprintf("%s，%s，%s %s %s %s %s", this.Name, this.Phone, this.Province, this.City, this.Areas, this.Addr)
}
