package models

import (
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/utils"
)

const (
	OrderStatusPrePay      = 0
	OrderStatusPaid        = 1
	OrderStatusFailed      = 2
	OrderStatusPreSend     = 3
	OrderStatusPreConfirm  = 4
	OrderStatusPreEvaluate = 5
	OrderStatusDone        = 6

	OrderTakeGoodTypeOnline  = 1
	OrderTakeGoodTypeOffline = 2
)

type Order struct {
	model.MongoModel `inline`
	OrderNo          string                 `json:"order_no" bson:"order_no" name:"订单号"`
	ItemCount        int64                  `json:"item_count" bson:"item_count" name:"订单商品数量"`
	OrderAmount      uint64                 `json:"order_amount" bson:"order_amount" name:"订单金额,单位分"`
	ActualAmount     uint64                 `json:"actual_amount" bson:"actual_amount" name:"实付金额,单位分"`
	OrderItems       []*OrderItem           `json:"order_items" bson:"order_items" name:"订单详情"`
	User             *AssociatedUser        `json:"user" bson:"user" name:"用户"`
	UserAddress      *AssociatedUserAddress `json:"user_address" bson:"user_address" name:"收货信息"`
	TakeGoodType     int                    `json:"take_good_type" bson:"take_good_type" name:"物流类型"`
	Logistics        []*Logistics           `json:"logistics" name:"物流信息"`
	Payment          *AssociatedPayment     `json:"payment" name:"支付信息"`
	Status           int                    `json:"status" name:"订单状态"`
}

func NewOrder() *Order {
	return &Order{}
}

func (this *Order) OriginName() string {
	return "订单出库"
}

func (this *Order) OriginModel() string {
	return utils.StructNameToSnakeAndPlural(this)
}

func (this *Order) OriginId() string {
	return this.GetID()
}

type OrderItem struct {
	Item  *AssociatedItem `json:"item"`
	Count int64           `json:"count"`
}

type Logistics struct {
	Enterprise string `json:"enterprise"`
	TrackNo    string `json:"track_no" bson:"track_no"`
}
