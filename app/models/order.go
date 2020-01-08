package models

import "go-shop-v2/pkg/db/model"

type Order struct {
	model.MongoModel
	OrderNo     string             `json:"order_no" bson:"order_no" name:"订单号"`
	OrderItems  []*OrderItem       `json:"order_items" bson:"order_items" name:"订单详情"`
	User        *AssociatedUser    `json:"user" name:"用户"`
	UserAddress *UserAddress       `json:"user_address" bson:"user_address" name:"收货信息"`
	Logistics   *Logistics         `json:"logistics" name:"物流信息"`
	Payment     *AssociatedPayment `json:"payment" name:"支付信息"`
	Status      int                `json:"status" name:"订单状态"`
}

type OrderItem struct {
	Item  *AssociatedItem `json:"item"`
	Count int             `json:"count"`
}

type Logistics struct {
	Enterprise string `json:"enterprise"`
	TrackNo    string `json:"track_no" bson:"track_no"`
}
