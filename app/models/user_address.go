package models

import "go-shop-v2/pkg/db/model"

type UserAddress struct {
	model.MongoModel `inline`
	UserId           string `json:"user_id" bson:"user_id"`
	ContactName      string `json:"contact_name" bson:"contact_name"`
	ContactPhone     string `json:"contact_phone" bson:"contact_phone"`
	Province         string `json:"province"`
	City             string `json:"city"`
	Areas            string `json:"areas"`
	Addr             string `json:"addr"`
	IsDefault        int    `json:"is_default" bson:"is_default"`
}

type AssociatedUserAddress struct {
	Id           string `json:"id"`
	ContactName  string `json:"contact_name" bson:"contact_name"`
	ContactPhone string `json:"contact_phone" bson:"contact_phone"`
	Province     string `json:"province"`
	City         string `json:"city"`
	Areas        string `json:"areas"`
	Addr         string `json:"addr"`
}

func (address *UserAddress) ToAssociated() *AssociatedUserAddress {
	return &AssociatedUserAddress{
		Id:           address.GetID(),
		ContactName:  address.ContactName,
		ContactPhone: address.ContactPhone,
		Province:     address.Province,
		City:         address.City,
		Areas:        address.Areas,
		Addr:         address.Addr,
	}
}
