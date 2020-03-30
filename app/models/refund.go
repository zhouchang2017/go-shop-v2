package models

import "go-shop-v2/pkg/db/model"

// 退款
type Refund struct {
	model.MongoModel `inline`
	OrderNo          string        `json:"order_no" bson:"order_no"`               // 订单号
	RefundOrderNo    string        `json:"refund_order_no" bson:"refund_order_no"` // 退款单号
	PaymentNo        string        `json:"payment_no" bson:"payment_no"`           // 微信订单号
	PaymentStatus    string        `json:"payment_status" bson:"payment_status"`   // 微信执行状态
	TotalAmount      uint64        `json:"total_amount" bson:"total_amount"`       // 退款金额
	Items            []*RefundItem `json:"items"`                                  // todo: 这里要加bson标签吗？
	// todo: add logistics no

}

type RefundItem struct {
	ItemId string `json:"item_id"` // 退款
	Qty    int64  `json:"qty"`     // 数量
	Amount uint64 `json:"amount"`  // 退款金额
}
