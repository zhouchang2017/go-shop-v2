package models

import "go-shop-v2/pkg/db/model"

const (
	RefundStatusPending    = iota // 等待后台处理
	RefundStatusApprove           // 后台处理通过
	RefundStatusProcessing        // 退款处理中
	RefundStatusFinished          // 退款完成
	RefundStatusRefused           // 拒绝
	RefundStatusCancel            // 取消
)

// 退款
type Refund struct {
	model.MongoModel `inline`
	OrderNo          string        `json:"order_no" bson:"order_no"`         // 订单号
	TotalAmount      uint64        `json:"total_amount" bson:"total_amount"` // 退款金额
	Items            []*RefundItem `json:"items"`
}

type RefundItem struct {
	ItemId string `json:"item_id"` // 退款
	Qty    string `json:"qty"`     // 数量
	Amount uint64 `json:"amount"`  // 退款金额
}
