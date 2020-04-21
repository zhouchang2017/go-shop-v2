package models

import (
	"fmt"
	"go-shop-v2/pkg/db/model"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
	"unicode/utf8"
)

const (
	OrderStatusFailed      = iota - 1 // 订单关闭
	OrderStatusPrePay                 // 等待付款
	OrderStatusPaid                   // 支付成功
	OrderStatusPreSend                // 等待发货
	OrderStatusPreConfirm             // 等待收货
	OrderStatusConfirm                // 确认收货
	OrderStatusPreEvaluate            // 待评价
	OrderStatusDone                   // 交易成功

	OrderTakeGoodTypeOnline  = 1
	OrderTakeGoodTypeOffline = 2
)

var OrderStatus []struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Class string `json:"class"`
	Type  string `json:"type"`
	Step  bool   `json:"step"`
} = []struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Class string `json:"class"`
	Type  string `json:"type"`
	Step  bool   `json:"step"`
}{
	{Name: "已取消", Value: OrderStatusFailed, Class: "bg-gray-400"},
	{Name: "待付款", Value: OrderStatusPrePay, Class: "bg-yellow-400", Type: "order", Step: true},
	{Name: "已付款", Value: OrderStatusPaid, Class: "bg-blue-400", Type: "order", Step: true},
	{Name: "待发货", Value: OrderStatusPreSend, Class: "bg-red-400", Type: "order", Step: true},
	{Name: "待收货", Value: OrderStatusPreConfirm, Class: "bg-green-200", Type: "order", Step: true},
	{Name: "已收货", Value: OrderStatusConfirm, Class: "bg-green-300", Type: "order", Step: true},
	{Name: "待评价", Value: OrderStatusPreEvaluate, Class: "bg-green-400", Type: "order", Step: true},
	{Name: "交易完成", Value: OrderStatusDone, Class: "bg-green-400", Type: "order", Step: true},
}

type Order struct {
	model.MongoModel `inline`
	OrderNo          string                 `json:"order_no" bson:"order_no" name:"订单号"`
	ItemCount        uint64                 `json:"item_count" bson:"item_count" name:"订单商品数量"`
	OrderAmount      uint64                 `json:"order_amount" bson:"order_amount" name:"订单金额,单位分"`
	ActualAmount     uint64                 `json:"actual_amount" bson:"actual_amount" name:"实付金额,单位分"`
	OrderItems       []*OrderItem           `json:"order_items" bson:"order_items" name:"订单详情"`
	User             *AssociatedUser        `json:"user" bson:"user" name:"用户"`
	UserAddress      *AssociatedUserAddress `json:"user_address" bson:"user_address" name:"收货信息"`
	TakeGoodType     int                    `json:"take_good_type" bson:"take_good_type" name:"物流类型"`
	Logistics        []*Logistics           `json:"logistics" name:"物流信息"`
	RefundMark       int                    `json:"refund_mark" bson:"refund_mark"`           // 退款标识 无退款 = 0、存在退款 = 1、完全退款 = 2
	RemainderAmount  uint64                 `json:"remainder_amount" bson:"remainder_amount"` // 退款后该笔订单剩余用户支付金额
	Refunds          []*Refund              `json:"refunds"`                                  // 退款信息
	Payment          *AssociatedPayment     `json:"payment" name:"支付信息"`
	Status           int                    `json:"status" name:"订单状态"`
	RefundChannel    bool                   `json:"refund_channel" bson:"refund_channel"`                 // 是否关闭退款通道
	CommentChannel   bool                   `json:"comment_channel" bson:"comment_channel"`               // 是否关闭评价通道
	PromotionInfo    *PromotionOverView     `json:"promotion_info" bson:"promotion_info"`                 // 促销信息
	ShipmentsAt      *time.Time             `json:"shipments_at" bson:"shipments_at"`                     // 发货时间
	CommentedAt      *time.Time             `json:"commented_at" bson:"commented_at"`                     // 评价时间
	CloseReason      *string                `json:"close_reason,omitempty" bson:"close_reason,omitempty"` // 订单取消原因
}

// 获取订单关联的product
func (o Order) GetProductIds() []string {
	ids := make([]string, 0)
	for _, item := range o.OrderItems {
		if exist, _ := utils.InArray(item.Item.Product.Id, ids); !exist {
			ids = append(ids, item.Item.Product.Id)
		}
	}
	return ids
}

// 订单是否可以退款
func (o Order) CanApplyRefund() error {
	if o.Status == OrderStatusPreSend && o.RefundChannel == true {
		return nil
	}
	if o.RefundChannel {
		return err2.Err422.F("当前订单不允许退款，退款通道已关闭")
	}
	return err2.Err422.F("当前订单不允许退款")
}

// 订单缩略图
func (o Order) GetAvatar() string {
	return o.OrderItems[0].Item.Avatar.Src()
}

// 订单第一件
type AggregateUnitLogistics struct {
	Items        *LogisticsItem `json:"items"`
	NoDelivery   bool           `json:"no_delivery" bson:"no_delivery"`     // 是否无需物流
	DeliveryName string         `json:"delivery_name" bson:"delivery_name"` // 物流公司名称
	DeliveryId   string         `json:"delivery_id" bson:"delivery_id"`     // 物流公司标识
	TrackNo      string         `json:"track_no" bson:"track_no"`           // 物流单号
	CreatedAt    time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" bson:"updated_at"`
}

type AggregateOrderItem struct {
	OrderId       primitive.ObjectID      `json:"order_id" bson:"order_id"`
	CreatedAt     time.Time               `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at" bson:"updated_at"`
	OrderNo       string                  `json:"order_no" bson:"order_no" name:"订单号"`
	ItemCount     int64                   `json:"item_count" bson:"item_count" name:"订单商品数量"`
	OrderAmount   uint64                  `json:"order_amount" bson:"order_amount" name:"订单金额,单位分"`
	ActualAmount  uint64                  `json:"actual_amount" bson:"actual_amount" name:"实付金额,单"`
	OrderItem     *OrderItem              `json:"order_item" bson:"order_item"`
	Logistics     *AggregateUnitLogistics `json:"logistics" name:"物流信息"`
	Payment       *AssociatedPayment      `json:"payment" name:"支付信息"`
	Status        int                     `json:"status" name:"订单状态"`
	PromotionInfo *PromotionOverView      `json:"promotion_info" bson:"promotion_info"` // 促销信息
	ShipmentsAt   *time.Time              `json:"shipments_at" bson:"shipments_at"`     // 发货时间
	CommentedAt   *time.Time              `json:"commented_at" bson:"commented_at"`     // 评价时间
}

func orderStatusText(status int) string {
	switch status {
	case OrderStatusFailed:
		return "已关闭"
	case OrderStatusPrePay:
		return "待付款"
	case OrderStatusPaid:
		return "支付成功"
	case OrderStatusPreSend:
		return "等待发货"
	case OrderStatusPreConfirm:
		return "已发货"
	case OrderStatusConfirm:
		return "已收货"
	case OrderStatusPreEvaluate:
		return "待评价"
	case OrderStatusDone:
		return "交易完成"
	}
	return "N/A"
}

// 订单状态
func (o Order) StatusText() string {
	return orderStatusText(o.Status)
}

// 设置关闭理由
func (o *Order) SetCloseReason(reason string) {
	if reason == "" {
		reason = "交易超时自动关闭"
	}
	o.CloseReason = &reason
}

// 获取订单支付金额
func (o Order) GetActualAmount() string {
	amount := float64(o.ActualAmount) / 100
	float := strconv.FormatFloat(amount, 'f', 2, 64)
	return fmt.Sprintf("￥%s", float)
}

// 获取订单商品名称
func (o Order) GoodsName(limit int) string {
	if len(o.OrderItems) > 1 {
		name := o.OrderItems[0].Item.Product.Name
		if limit == -1 {
			return fmt.Sprintf("%s(等商品)", name)
		}
		if utf8.RuneCountInString(name) > limit-8 {
			subString := utils.SubString(name, 0, limit-8)
			return fmt.Sprintf("%s...(等商品)", subString)
		}
		return fmt.Sprintf("%s(等商品)", name)
	}
	name := o.OrderItems[0].Item.Product.Name
	if limit == -1 {
		return fmt.Sprintf("%s", name)
	}
	if utf8.RuneCountInString(name) > limit-3 {
		subString := utils.SubString(name, 0, limit-3)
		return fmt.Sprintf("%s...", subString)
	}
	return fmt.Sprintf("%s", name)
}

// 订单是否已关闭
func (o *Order) StatusIsFailed() bool {
	return o.Status == OrderStatusFailed
}

// 状态设置为取消
func (o *Order) StatusToFailed() error {
	if o.Status == OrderStatusPrePay {
		o.Status = OrderStatusFailed
		return nil
	}
	return err2.Err422.F("当前订单状态[%d]不允许取消", o.Status)
}

// 判断是否能评论
func (o Order) CanComment() bool {
	if (o.Status == OrderStatusPreEvaluate || o.Status == OrderStatusConfirm) && o.CommentedAt == nil {
		return true
	}
	return false
}

// 订单总计商品数量
func (o Order) ItemsQty() (count uint64) {
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

func (this *Order) FindItem(id string) *OrderItem {
	for _, item := range this.OrderItems {
		if item.Item.Id == id {
			return item
		}
	}
	return nil
}

type OrderItem struct {
	Item            *AssociatedItem    `json:"item"`
	Count           uint64             `json:"count"`                                    // 购买数量
	Price           uint64             `json:"price"`                                    // item单品优惠价格，受Promotion.Type = 0 的影响
	TotalAmount     uint64             `json:"total_amount" bson:"total_amount"`         // 实际支付价格
	PromotionInfo   *ItemPromotionInfo `json:"promotion_info" bson:"promotion_info"`     // 冗余促销信息
	RemainderAmount uint64             `json:"remainder_amount" bson:"remainder_amount"` // 该订单商品记录剩余用户支付金额
	RemainderQty    uint64             `json:"remainder_qty" bson:"remainder_qty"`       // 剩余数量 Quantity
	Refunding       bool               `json:"refunding"`                                // 是否存在进行中的退款
}

type OrderCountByStatus struct {
	Status int   `json:"status" bson:"_id"`
	Count  int64 `json:"count"`
}
