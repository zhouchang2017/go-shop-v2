package models

import (
	"fmt"
	"go-shop-v2/pkg/db/model"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/utils"
	"time"
)

const (
	OrderStatusFailed      = iota - 1 // 订单关闭
	OrderStatusPrePay                 // 等待付款
	OrderStatusPaid                   // 支付成功
	OrderStatusPreSend                // 等待发货
	OrderStatusPartSend               // 部分发货
	OrderStatusPreConfirm             // 等待收货
	OrderStatusPreEvaluate            // 待评价
	OrderStatusDone                   // 交易完成

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
	ShipmentsAt      *time.Time             `json:"shipments_at" bson:"shipments_at"`     // 发货时间
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

func (this *Order) findItem(id string) *OrderItem {
	for _, item := range this.OrderItems {
		if item.Item.Id == id {
			return item
		}
	}
	return nil
}

// 计算所有包裹中该商品的数量
func (this *Order) countItem(id string) (count int64) {
	for _, item := range this.Logistics {
		if i := item.findItem(id); i != nil {
			count += i.Count
		}
	}
	return count
}

// 计算物流状态
func (this *Order) refreshShipmentStatus() {
	totalItemCount := this.ItemsQty()
	var totalShipmentCount int64
	for _, item := range this.Logistics {
		totalShipmentCount += item.itemsQty()
	}
	if totalShipmentCount == 0 {
		// 未发货
		this.Status = OrderStatusPreSend
		return
	}
	if totalItemCount == totalShipmentCount {
		// 全部发货完成
		this.Status = OrderStatusPreConfirm
		return
	}
	if totalShipmentCount > 0 && totalShipmentCount < totalItemCount {
		// 部分发货
		this.Status = OrderStatusPartSend
	}
	return
}

// 发货
func (this *Order) Shipment(opts ...*LogisticsOption) error {
	for _,opt:=range opts {
		if err := opt.isValid(); err != nil {
			return err
		}
		item := this.findItem(opt.ItemId)
		if item == nil {
			return err2.Err422.F("该订单中不存在itemId[%s]", opt.ItemId)
		}
		// 验证产品数量
		// 已发货数量
		countItem := this.countItem(opt.ItemId)
		if opt.Count > item.Count-countItem {
			return err2.Err422.F("产品数量超出该订单购买商品数量，发货数量[%d],订单中该商品数量[%d],以发货数量", opt.Count, item.Count, countItem)
		}
		var logistics *Logistics
		for _, l := range this.Logistics {
			if l.Enterprise == opt.Enterprise && l.TrackNo == opt.TrackNo {
				logistics = l
				break
			}
		}
		if logistics != nil {
			if err := logistics.addItem(opt.ItemId, opt.Count, opt.ShopId); err != nil {
				return err
			}
		} else {
			logistics = &Logistics{
				Enterprise: opt.Enterprise,
				TrackNo:    opt.TrackNo,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			if err := logistics.addItem(opt.ItemId, opt.Count, opt.ShopId); err != nil {
				return err
			}
			this.Logistics = append(this.Logistics, logistics)
		}
	}

	this.refreshShipmentStatus()
	return nil
}

type OrderItem struct {
	Item          *AssociatedItem    `json:"item"`
	Count         int64              `json:"count"`                                // 购买数量
	Price         int64              `json:"price"`                                // item单品优惠价格，受Promotion.Type = 0 的影响
	Amount        int64              `json:"amount"`                               // 实际支付价格
	PromotionInfo *ItemPromotionInfo `json:"promotion_info" bson:"promotion_info"` // 冗余促销信息
}

// 发货选项结构
type LogisticsOption struct {
	Enterprise string `json:"enterprise"`
	TrackNo    string `json:"track_no"`
	ItemId     string `json:"item_id"`
	Count      int64  `json:"count"`
	ShopId     string `json:"shop_id"`
}

func (l LogisticsOption) isValid() error {
	if l.ItemId == "" {
		return err2.Err422.F("缺少发货商品id")
	}
	if l.Count == 0 {
		return err2.Err422.F("发货商品数量必须大于0")
	}
	if l.ShopId == "" {
		return err2.Err422.F("缺少寄件方")
	}
	if l.Enterprise == "" {
		return err2.Err422.F("缺少物流公司")
	}
	if l.TrackNo == "" {
		return err2.Err422.F("缺少物流单号")
	}
	return nil
}

// 包裹明细
type LogisticsItem struct {
	ItemId string `json:"item_id" bson:"item_id"` // 商品id
	Count  int64  `json:"count"`                  // 数量
	ShopId string `json:"shop_id" bson:"shop_id"` // 出货门店
}

// 物流
type Logistics struct {
	Items      []*LogisticsItem `json:"items"`
	Enterprise string           `json:"enterprise"`
	TrackNo    string           `json:"track_no" bson:"track_no"`
	CreatedAt  time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at" bson:"updated_at"`
}

// 计算物品总数
func (l Logistics) itemsQty() (count int64) {
	for _, item := range l.Items {
		count += item.Count
	}
	return count
}

// 搜索包裹产品
func (l *Logistics) findItem(id string) *LogisticsItem {
	for _, item := range l.Items {
		if item.ItemId == id {
			return item
		}
	}
	return nil
}

// 添加物品到包裹内
func (l *Logistics) addItem(itemId string, count int64, shopId string) error {
	for _, item := range l.Items {
		if item.ShopId != shopId {
			return fmt.Errorf("单个包裹寄件方只允许存在1个,当前包裹寄件方为[shopId=%s],添加物品寄件方为[%s]", item.ShopId, shopId)
		}
		if item.ItemId == itemId {
			item.Count += count
			l.UpdatedAt = time.Now()
			return nil
		}
	}
	l.Items = append(l.Items, &LogisticsItem{
		ItemId: itemId,
		Count:  count,
		ShopId: shopId,
	})
	l.UpdatedAt = time.Now()
	return nil
}
