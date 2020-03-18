package models

import (
	"go-shop-v2/pkg/db/model"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/utils"
)

const (
	OrderStatusPrePay      = 0 // 等待付款
	OrderStatusPaid        = 1 // 支付成功
	OrderStatusFailed      = 2 // 交易失败
	OrderStatusPreSend     = 3 // 等待发货
	OrderStatusPreConfirm  = 4 // 等待收货
	OrderStatusPreEvaluate = 5 // 待评价
	OrderStatusDone        = 6 // 交易完成

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
	PromotionInfo    *PromotionOverView     `json:"promotion_info" bson:"promotion_info"` // 促销信息
}

// 状态设置为取消
func (o *Order) StatusToFailed() error {
	if o.Status == OrderStatusPrePay {
		o.Status = OrderStatusFailed
		return nil
	}
	return err2.Err422.F("当前订单状态[%d]不允许取消", o.Status)
}

// 订单总计商品数量
func (o Order) ItemsQty() (count int64) {
	for _, item := range o.OrderItems {
		count += item.Count
	}
	return count
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
	Item          *AssociatedItem    `json:"item"`
	Count         int64              `json:"count"`                                // 购买数量
	Price         int64              `json:"price"`                                // item单品优惠价格，受Promotion.Type = 0 的影响
	Amount        int64              `json:"amount"`                               // 实际支付价格
	PromotionInfo *ItemPromotionInfo `json:"promotion_info" bson:"promotion_info"` // 冗余促销信息
}

type Logistics struct {
	Enterprise string `json:"enterprise"`
	TrackNo    string `json:"track_no" bson:"track_no"`
}
