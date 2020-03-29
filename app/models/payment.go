package models

import (
	"bytes"
	"go-shop-v2/pkg/db/model"
	"strconv"
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

func (p *Payment) SetPaymentAt(time2 time.Time) *Payment {
	p.PaymentAt = &time2
	return p
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
	Amount    uint64     `json:"amount" name:"金额，单位分"`
	PaymentNo string     `json:"payment_no" bson:"payment_no" name:"支付单号"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"` // 创建时间
	PaymentAt *time.Time `json:"payment_at" bson:"payment_at"` // 支付时间
}

// 日支付统计
type DayPaymentCount struct {
	Date        string `json:"date,omitempty" bson:"_id,omitempty"`
	TotalAmount Amount `json:"total_amount" bson:"total_amount"` // 总支付金额
	Count       int64  `json:"count" bson:"count"`               // 支付笔数
}

type Amount int64

func (a Amount) MarshalJSON() ([]byte, error) {
	amount := float64(a) / 100
	price := strconv.FormatFloat(amount, 'f', 2, 64)
	bufferString := bytes.NewBufferString(`"`)
	bufferString.WriteString(price)
	bufferString.WriteString(`"`)
	return bufferString.Bytes(), nil
}
