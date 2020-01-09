package models

import "go-shop-v2/pkg/db/model"

type UserAddress struct {
	model.MongoModel `inline`
	UserId           string `json:"user_id" bson:"user_id"`
	ContactName      string `json:"contact_name" bson:"contact_name"`
	ContactPhone     string `json:"contact_phone" bson:"contact_phone"`
	Country          string `json:"country"`
	Province         string `json:"province"`
	City             string `json:"city"`
	Address          string `json:"address"`
	IsDefault        int    `json:"is_default" bson:"is_default"`
}
