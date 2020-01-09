package models

import "go-shop-v2/pkg/db/model"

type Payment struct {
	model.MongoModel
	OrderNo        string `json:"order_no" bson:"order_no"`
	Platform       string `json:"platform" name:"支付平台:微信/支付宝"`
	Title          string `json:"title" name:"购买的商品标题"`
	Amount         int    `json:"amount" name:"金额，单位分"`
	ExtendedUserId string `json:"extended_user_id" bson:"extended_user_id" name:"平台用户标识"`
	PrePaymentNo   string `json:"pre_payment_no" bson:"pre_payment_no" name:"预下单号"`
	PaymentNo      string `json:"payment_no" bson:"payment_no" name:"支付单号"`
}

type AssociatedPayment struct {
	Platform  string `json:"platform" name:"支付平台:微信/支付宝"`
	Amount    int    `json:"amount" name:"金额，单位分"`
	PaymentNo string `json:"payment_no" bson:"payment_no" name:"支付单号"`
}
