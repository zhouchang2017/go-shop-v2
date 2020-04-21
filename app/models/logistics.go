package models

import (
	"fmt"
	err2 "go-shop-v2/pkg/err"
	"reflect"
	"time"
)

// 发货选项结构
type LogisticsOption struct {
	NoDelivery  bool               `json:"no_delivery" form:"no_delivery"` // 无需物流
	DeliveryId  string             `json:"delivery_id" form:"delivery_id"` // 物流公司编号
	TrackNo     string             `json:"track_no" form:"track_no"`       // 物流单号
	ShopId      string             `json:"shop_id" form:"shop_id"`
	Items       []*LogisticsItem   `json:"items"`
	WaybillID   string             `json:"-" form:"-"` // 小程序物流助手，运单ID，下单成功时返回
	WaybillData []*WaybillDataItem `json:"-" form:"-"` // 小程序物流助手，运单信息，下单成功时返回
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
	if l.WaybillID != "" {
		return nil
	}
	if l.TrackNo == "" {
		return err2.Err422.F("缺少物流单号")
	}

	return nil
}

// 包裹明细
type LogisticsItem struct {
	ItemId string `json:"item_id" form:"item_id" bson:"item_id"` // 商品id
	Count  uint64  `json:"count"`                                 // 数量
	ShopId string `json:"shop_id" form:"shop_id" bson:"shop_id"` // 出货门店
}

func (l LogisticsItem) equal(item LogisticsItem) bool {
	if l.ItemId == item.ItemId && l.Count == item.Count && l.ShopId == item.ShopId {
		return true
	}
	return false
}

var LogisticsInfos = []LogisticsInfo{
	{DeliveryName: "百世快递", DeliveryId: "BEST", Services: []LogisticsInfoService{
		{Type: 1, Name: "普通快递"},
	}},
	{DeliveryName: "中国邮政速递物流(EMS)", DeliveryId: "EMS", Services: []LogisticsInfoService{
		{Type: 6, Name: "标准快递"},
		{Type: 9, Name: "快递包裹"},
	}},
	{DeliveryName: "品骏物流", DeliveryId: "PJ", Services: []LogisticsInfoService{
		{Type: 1, Name: "普通快递"},
	}},
	{DeliveryName: "顺丰速运", DeliveryId: "SF", Services: []LogisticsInfoService{
		{Type: 0, Name: "标准快递"},
		{Type: 1, Name: "顺丰即日"},
		{Type: 2, Name: "顺丰次晨"},
		{Type: 3, Name: "顺丰标快"},
		{Type: 4, Name: "顺丰特惠"},
	}},
	{DeliveryName: "圆通速递", DeliveryId: "YTO", Services: []LogisticsInfoService{
		{Type: 0, Name: "普通快递"},
		{Type: 1, Name: "圆准达"},
	}},
	{DeliveryName: "韵达快递", DeliveryId: "YUNDA", Services: []LogisticsInfoService{
		{Type: 0, Name: "标准快件"},
	}},
	{DeliveryName: "中通快递", DeliveryId: "ZTO", Services: []LogisticsInfoService{
		{Type: 0, Name: "标准快件"},
	}},
	{DeliveryName: "申通快递", DeliveryId: "STO", Services: []LogisticsInfoService{
		{Type: 1, Name: "标准快递"},
	}},
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
	DeliveryName string                 `json:"delivery_name"` // 快递公司
	DeliveryId   string                 `json:"delivery_id"`   // 快递公司id
	BizID        string                 `json:"biz_id"`
	Services     []LogisticsInfoService `json:"services"`
}

type LogisticsInfoService struct {
	Type uint8  `json:"type"`
	Name string `json:"name"`
}

// 物流
type Logistics struct {
	Items        []*LogisticsItem   `json:"items"`
	NoDelivery   bool               `json:"no_delivery" bson:"no_delivery"`                   // 是否无需物流
	DeliveryName string             `json:"delivery_name" bson:"delivery_name"`               // 物流公司名称
	DeliveryId   string             `json:"delivery_id" bson:"delivery_id"`                   // 物流公司标识
	TrackNo      string             `json:"track_no" bson:"track_no"`                         // 物流单号
	WaybillID    string             `json:"waybill_id,omitempty" bson:"waybill_id,omitempty"` // 小程序物流助手下单单号
	WaybillData  []*WaybillDataItem `json:"waybill_data,omitempty" bson:"waybill_data,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type WaybillDataItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (l Logistics) Equal(logistics Logistics) bool {
	if l.NoDelivery == logistics.NoDelivery && l.DeliveryId == logistics.DeliveryId && l.TrackNo == logistics.TrackNo {
		return reflect.DeepEqual(l.Items, logistics.Items)
	}
	return false
}

// 计算物品总数
func (l Logistics) itemsQty() (count uint64) {
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
func (l *Logistics) addItem(itemId string, count uint64, shopId string) error {
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

// 计算物流状态
func (this *Order) refreshShipmentStatus() {
	totalItemCount := this.ItemsQty()
	var totalShipmentCount uint64
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
			orderTotal uint64
			count      uint64
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
			if opt.WaybillID != "" {
				logistics.TrackNo = opt.WaybillID
				logistics.WaybillID = opt.WaybillID
				logistics.WaybillData = opt.WaybillData
			}
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

// 删除物流
func (this *Order) RemoveShipment(deliveryId string, trackNo string) {
	// 移除该订单物流
	var index int
	for ind, logistic := range this.Logistics {
		if logistic.DeliveryId == deliveryId && logistic.TrackNo == trackNo {
			index = ind
			break
		}
	}
	this.Logistics = append(this.Logistics[:index], this.Logistics[index+1:]...)
	this.refreshShipmentStatus()
}