package models

import (
	"go-shop-v2/pkg/db/model"
	"time"
)

type Payment struct {
	model.MongoModel `inline`
	OrderNo          string     `json:"order_no" bson:"order_no"`
	Platform         string     `json:"platform" name:"支付平台:微信/支付宝"`
	Title            string     `json:"title" name:"购买的商品标题"`
	Amount           uint64     `json:"amount" name:"金额，单位分"`
	ExtendedUserId   string     `json:"extended_user_id" bson:"extended_user_id" name:"平台用户标识"`
	PrePaymentNo     string     `json:"pre_payment_no" bson:"pre_payment_no" name:"预下单号"`
	PaymentNo        string     `json:"payment_no" bson:"payment_no" name:"支付单号"`
	PaymentAt        *time.Time `json:"payment_at" bson:"payment_at"` // 支付时间
}

func (p Payment) ToAssociated() *AssociatedPayment {
	return &AssociatedPayment{
		Platform:  p.Platform,
		Amount:    p.Amount,
		PaymentNo: p.PaymentNo,
		CreatedAt: p.CreatedAt,
		PaymentAt: p.PaymentAt,
	}
}

type AssociatedPayment struct {
	Platform  string     `json:"platform" name:"支付平台:微信/支付宝"`
	Amount    uint64        `json:"amount" name:"金额，单位分"`
	PaymentNo string     `json:"payment_no" bson:"payment_no" name:"支付单号"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"` // 创建时间
	PaymentAt *time.Time `json:"payment_at" bson:"payment_at"` // 支付时间
}
