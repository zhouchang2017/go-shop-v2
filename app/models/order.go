package models

import (
	"fmt"
	"go-shop-v2/pkg/db/model"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/utils"
	"reflect"
	"time"
)

const (
	OrderStatusFailed      = iota - 1 // 订单关闭
	OrderStatusPrePay                 // 等待付款
	OrderStatusPaid                   // 支付成功
	OrderStatusPreSend                // 等待发货
	OrderStatusPreConfirm             // 等待收货
	OrderStatusDone                   // 交易完成
	OrderStatusPreEvaluate            // 待评价

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
	CommentedAt      *time.Time             `json:"commented_at" bson:"commented_at"`     // 评价时间
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
	if (o.Status == OrderStatusPreEvaluate || o.Status == OrderStatusDone) && o.CommentedAt == nil {
		return true
	}
	return false
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

func (this *Order) FindItem(id string) *OrderItem {
	for _, item := range this.OrderItems {
		if item.Item.Id == id {
			return item
		}
	}
	return nil
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
		this.Status = OrderStatusPreSend
	}
	return
}

// 发货
func (this *Order) Shipment(opts ...*LogisticsOption) error {
	var shipments []*Logistics
	for _, opt := range opts {
		if err := opt.isValid(); err != nil {
			return err
		}
		type itemCountMapItem struct {
			orderTotal int64
			count      int64
		}
		var itemCountMap = map[string]*itemCountMapItem{}
		for _, i := range opt.Items {
			item := this.FindItem(i.ItemId)
			if item == nil {
				return err2.Err422.F("该订单中不存在itemId[%s]", i.ItemId)
			}
			if itemCountMap[i.ItemId] != nil {
				itemCountMap[i.ItemId].count += i.Count
			} else {
				itemCountMap[i.ItemId] = &itemCountMapItem{
					orderTotal: item.Count,
					count:      i.Count,
				}
			}
		}
		for key, value := range itemCountMap {
			if value.orderTotal-value.count < 0 {
				return err2.Err422.F("itemId[%s]发货数量溢出，总发货数量%d，实际发货数量%d\n", key, value.orderTotal, value.count)
			}
		}
		var logistics *Logistics

		logistics = &Logistics{
			NoDelivery: opt.NoDelivery,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if !opt.NoDelivery {
			info := FindLogisticsInfo(opt.DeliveryId)
			logistics.DeliveryId = opt.DeliveryId
			logistics.DeliveryName = info.DeliveryName
			logistics.TrackNo = opt.TrackNo
		}

		for _, i := range opt.Items {
			if err := logistics.addItem(i.ItemId, i.Count, opt.ShopId); err != nil {
				return err
			}
		}

		shipments = append(shipments, logistics)

		//for index, l := range this.Logistics {
		//	if index+1 > len(opts) {
		//		break
		//	}
		//	mockLogistics := Logistics{
		//		Items:      opt.Items,
		//		NoDelivery: opt.NoDelivery,
		//		DeliveryId: opt.DeliveryId,
		//		TrackNo:    opt.TrackNo,
		//	}
		//	if l.DeliveryId == opt.DeliveryId && l.TrackNo == opt.TrackNo && l.NoDelivery == opt.NoDelivery {
		//		// 单号未发生变化情况
		//		if l.Equal(mockLogistics) {
		//			// 无更新
		//			logistics = l
		//		} else {
		//			// 商品存在更新
		//			// 先清空包裹
		//			l.Items = []*LogisticsItem{}
		//			// 添加商品到包裹
		//			for _, i := range opt.Items {
		//				if err := l.addItem(i.ItemId, i.Count, opt.ShopId); err != nil {
		//					return err
		//				}
		//			}
		//			l.UpdatedAt = time.Now()
		//		}
		//		shipments = append(shipments, l)
		//		continue
		//	}
		//	if l.Equal(mockLogistics) {
		//		//  更新了物流信息
		//		logistics = l
		//		logistics.NoDelivery = opt.NoDelivery
		//		if !logistics.NoDelivery {
		//			logistics.TrackNo = opt.TrackNo
		//			logistics.DeliveryId = opt.DeliveryId
		//			logistics.DeliveryName = FindLogisticsInfo(opt.DeliveryId).DeliveryName
		//		} else {
		//			logistics.TrackNo = ""
		//			logistics.DeliveryId = ""
		//			logistics.DeliveryName = ""
		//		}
		//		logistics.UpdatedAt = time.Now()
		//		shipments = append(shipments, logistics)
		//		continue
		//	}
		//}
		//if logistics == nil {
		//
		//	logistics = &Logistics{
		//		NoDelivery: opt.NoDelivery,
		//		CreatedAt:  time.Now(),
		//		UpdatedAt:  time.Now(),
		//	}
		//	if !opt.NoDelivery {
		//		info := FindLogisticsInfo(opt.DeliveryId)
		//		logistics.DeliveryId = opt.DeliveryId
		//		logistics.DeliveryName = info.DeliveryName
		//		logistics.TrackNo = opt.TrackNo
		//	}
		//
		//	for _, i := range opt.Items {
		//		if err := logistics.addItem(i.ItemId, i.Count, opt.ShopId); err != nil {
		//			return err
		//		}
		//	}
		//
		//	shipments = append(shipments, logistics)
		//}
	}
	this.Logistics = shipments
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
	NoDelivery bool             `json:"no_delivery" form:"no_delivery"` // 无需物流
	DeliveryId string           `json:"delivery_id" form:"delivery_id"` // 物流公司编号
	TrackNo    string           `json:"track_no" form:"track_no"`       // 物流单号
	ShopId     string           `json:"shop_id" form:"shop_id"`
	Items      []*LogisticsItem `json:"items"`
}

func (l LogisticsOption) isValid() error {

	if len(l.Items) == 0 {
		return err2.Err422.F("缺少发货商品")
	}

	for _, item := range l.Items {
		if item.ItemId == "" {
			return err2.Err422.F("缺少发货商品id")
		}
		if item.Count <= 0 {
			return err2.Err422.F("发货商品数量必须大于0")
		}
	}

	if l.ShopId == "" {
		return err2.Err422.F("缺少寄件方")
	}
	if l.NoDelivery {
		return nil
	}
	if l.DeliveryId == "" {
		return err2.Err422.F("缺少物流公司")
	}
	if l.TrackNo == "" {
		return err2.Err422.F("缺少物流单号")
	}
	return nil
}

// 包裹明细
type LogisticsItem struct {
	ItemId string `json:"item_id" form:"item_id" bson:"item_id"` // 商品id
	Count  int64  `json:"count"`                                 // 数量
	ShopId string `json:"shop_id" form:"shop_id" bson:"shop_id"` // 出货门店
}

func (l LogisticsItem) equal(item LogisticsItem) bool {
	if l.ItemId == item.ItemId && l.Count == item.Count && l.ShopId == item.ShopId {
		return true
	}
	return false
}

var LogisticsInfos = []LogisticsInfo{
	{"百世快递", "BEST"},
	{"中国邮政速递物流(EMS)", "EMS"},
	{"品骏物流", "PJ"},
	{"顺丰速运", "SF"},
	{"圆通速递", "YTO"},
	{"韵达快递", "YUNDA"},
	{"中通快递", "ZTO"},
	{"申通快递", "STO"},
}

func FindLogisticsInfo(id string) LogisticsInfo {
	for _, item := range LogisticsInfos {
		if item.DeliveryId == id {
			return item
		}
	}
	return LogisticsInfo{DeliveryName: "未知", DeliveryId: "Unknown"}
}

type LogisticsInfo struct {
	DeliveryName string `json:"delivery_name"` // 快递公司
	DeliveryId   string `json:"delivery_id"`   // 快递公司id
}

// 物流
type Logistics struct {
	Items        []*LogisticsItem `json:"items"`
	NoDelivery   bool             `json:"no_delivery" bson:"no_delivery"`     // 是否无需物流
	DeliveryName string           `json:"delivery_name" bson:"delivery_name"` // 物流公司名称
	DeliveryId   string           `json:"delivery_id" bson:"delivery_id"`     // 物流公司标识
	TrackNo      string           `json:"track_no" bson:"track_no"`           // 物流单号
	CreatedAt    time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at" bson:"updated_at"`
}

func (l Logistics) Equal(logistics Logistics) bool {
	if l.NoDelivery == logistics.NoDelivery && l.DeliveryId == logistics.DeliveryId && l.TrackNo == logistics.TrackNo {
		return reflect.DeepEqual(l.Items, logistics.Items)
	}
	return false
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
