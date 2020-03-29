package models

import (
	"errors"
	"fmt"
	"go-shop-v2/app/lbs"
	"go-shop-v2/pkg/db/model"
)

type UserAddress struct {
	model.MongoModel `inline`
	UserId           string `json:"user_id" bson:"user_id"`
	ContactName      string `json:"contact_name" bson:"contact_name"`
	ContactPhone     string `json:"contact_phone" bson:"contact_phone"`
	Province         string `json:"province"`
	City             string `json:"city"`
	Areas            string `json:"areas"`
	Addr             string `json:"addr"`
	IsDefault        bool   `json:"is_default" bson:"is_default"`
}

func (this UserAddress) String() string {
	return fmt.Sprintf("%s %s %s%s%s%s", this.ContactName, this.ContactPhone, this.Province, this.City, this.Areas, this.Addr)
}

func (this UserAddress) AddressString() string {
	return fmt.Sprintf("%s%s%s%s", this.Province, this.City, this.Areas, this.Addr)
}

func (this UserAddress) Location() (loc *Location, err error) {
	loc = &Location{}
	if lbs.SDK != nil {
		location, err := lbs.SDK.ParseAddress(this.AddressString())
		if err != nil {
			return nil, err
		}
		loc.Lat = location.Lat
		loc.Lng = location.Lng
		return loc, nil
	}
	return nil, errors.New("lbs skd is nil")
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

func (this AssociatedUserAddress) String() string {
	return fmt.Sprintf("%s %s %s%s%s%s", this.ContactName, this.ContactPhone, this.Province, this.City, this.Areas, this.Addr)
}

func (this AssociatedUserAddress) AddressString() string {
	return fmt.Sprintf("%s%s%s%s", this.Province, this.City, this.Areas, this.Addr)
}

func (this AssociatedUserAddress) Location() (loc *Location, err error) {
	loc = &Location{}
	if lbs.SDK != nil {
		location, err := lbs.SDK.ParseAddress(this.AddressString())
		if err != nil {
			return nil, err
		}
		loc.Lat = location.Lat
		loc.Lng = location.Lng
		return loc, nil
	}
	return nil, errors.New("lbs skd is nil")
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
