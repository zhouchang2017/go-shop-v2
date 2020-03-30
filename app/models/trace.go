package models

import (
	"go-shop-v2/pkg/db/model"
	"time"
)

// 物流追踪
type Trace struct {
	model.MongoModel `inline`
	OrderNo          string         `json:"order_no" bson:"order_no"`
	ToUserName       string         `json:"to_user_name" bson:"to_user_name"`
	FromUserName     string         `json:"from_user_name" bson:"from_user_name"`
	CreateTime       time.Time      `json:"create_time" bson:"create_time"`
	MsgType          string         `json:"msg_type" bson:"msg_type"`
	Event            string         `json:"event"`
	DeliveryID       string         `json:"delivery_id" bson:"delivery_id"`
	WayBillId        string         `json:"way_bill_id" bson:"way_bill_id"`
	Version          int            `json:"version"`
	Count            int            `json:"count"`
	Actions          []*TrackAction `json:"actions"`
}
type TrackAction struct {
	ActionTime time.Time `json:"action_time" bson:"action_time"`
	ActionType string    `json:"action_type" bson:"action_type"`
	ActionMsg  string    `json:"action_msg" bson:"action_msg"`
}

//type ActionType string
//
//func (a ActionType) MarshalJSON() ([]byte, error) {
//	if string(i) == "" {
//		return bytes.NewBufferString("null").Bytes(), nil
//	}
//	bufferString := bytes.NewBufferString(`"`)
//	bufferString.WriteString(i.Src())
//	bufferString.WriteString(`"`)
//	return bufferString.Bytes(), nil
//}
//
