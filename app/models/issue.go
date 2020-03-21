package models

import "go-shop-v2/pkg/qiniu"

const (
	IssueStatusOpen = iota + 1
	IssueStatusResolve
	IssueStatusUnResolve
	IssueStatusClosed
)

// 工单
type Issue struct {
	OrderId string         `json:"order_id" bson:"order_id"`
	Title   string         `json:"title"`   // 问题
	Content string         `json:"content"` // 问题描述
	Images  []*qiniu.Image `json:"images"`  // 图片，限制5张
	Status  int            `json:"status"`
}

type IssueItem struct {
	User    *IssueUser     `json:"user"`
	Content string         `json:"content"`
	Images  []*qiniu.Image `json:"images"`
}

type IssueUser struct {
	Type   string `json:"type"`
	UserId string `json:"user_id" bson:"user_id"`
}
