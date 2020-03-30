package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/medivhzhan/weapp/v2"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/cache/redis"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/wechat"
	"time"
)

// 微信小程序物流服务
type LogisticsService struct {
	orderRep *repositories.OrderRep
	traceRep *repositories.TraceRep
}

func NewLogisticsService(orderRep *repositories.OrderRep, traceRep *repositories.TraceRep) *LogisticsService {
	return &LogisticsService{orderRep: orderRep, traceRep: traceRep}
}

const logisticsCacheKey = "go-shop-my-delivery"

// 小程序物流助手，已绑定物流公司列表
func (this *LogisticsService) GetAllDelivery() (res []*models.LogisticsInfo, err error) {
	res = make([]*models.LogisticsInfo, 0)
	if wechat.SDK == nil {
		return res, nil
	}
	if redis.GetConFn() != nil {
		result, err := redis.GetConFn().Get(logisticsCacheKey).Result()
		if err == nil {
			if result != "" {
				if err := json.Unmarshal([]byte(result), &res); err == nil {
					return res, nil
				}
			}
		}
	}
	if account, err := wechat.SDK.GetAllAccount(); err == nil {
		infos := make([]*models.LogisticsInfo, 0)
		delivery, err := wechat.SDK.GetAllDelivery()
		if err != nil {
			for _, item := range models.LogisticsInfos {
				infos = append(infos, &models.LogisticsInfo{
					DeliveryName: item.DeliveryName,
					DeliveryId:   item.DeliveryId,
					BizID:        "",
					Services:     item.Services,
				})
			}
		} else {
			for _, item := range delivery.Data {

				info := &models.LogisticsInfo{
					DeliveryName: item.Name,
					DeliveryId:   item.ID,
				}

				for _, i := range models.LogisticsInfos {
					if i.DeliveryId == info.DeliveryId {
						info.Services = i.Services
						info.BizID = i.BizID
					}
				}

				infos = append(infos, info)
			}
		}

		for _, data := range account.List {
			for _, info := range infos {
				if data.DeliveryID == info.DeliveryId && data.StatusCode == 0 {
					res = append(res, &models.LogisticsInfo{
						DeliveryName: info.DeliveryName,
						DeliveryId:   data.DeliveryID,
						BizID:        data.BizID,
						Services:     info.Services,
					})
				}
			}
		}

		if redis.GetConFn() != nil {
			if marshal, err := json.Marshal(res); err == nil {
				redis.GetConFn().Set(logisticsCacheKey, marshal, time.Minute*60)
			}
		}
	}
	return res, nil
}

type CreateExpressOrderOption struct {
	OrderId             string                          `json:"order_id" form:"order_id" binding:"required"`
	DeliveryId          string                          `json:"delivery_id" form:"delivery_id" binding:"required"`
	DeliveryServiceType uint8                           `json:"delivery_service_type" form:"delivery_service_type"`
	CustomRemark        string                          `json:"custom_remark" form:"custom_remark"`
	Shop                *models.Shop                    `json:"shop" form:"shop" binding:"required"`
	Items               []*createExpressOrderItemOption `json:"items" binding:"required"`
	Weight              float64                         `json:"weight"`
	SpaceX              float64                         `json:"space_x" form:"space_x"`
	SpaceY              float64                         `json:"space_y" form:"space_y"`
	SpaceZ              float64                         `json:"space_z" form:"space_z"`
	UseInsured          int                             `json:"use_insured" form:"use_insured"`     // 是否保价
	InsuredValue        int64                           `json:"insured_value" form:"insured_value"` // 保价金额
	ExpectTime          time.Time                       `json:"expect_time" form:"expect_time"`     // 上门取件时间
}

type createExpressOrderItemOption struct {
	ItemId string `json:"item_id" form:"item_id" binding:"required"`
	Count  int64  `json:"count" form:"count" binding:"required,min=1"`
}

func (c CreateExpressOrderOption) IsValid() error {
	if c.OrderId == "" {
		return err2.Err422.F("缺少order id")
	}
	if c.DeliveryId == "nil" {
		return err2.Err422.F("缺少delivery id")
	}
	if c.Shop == nil {
		return err2.Err422.F("缺少发货门店")
	}
	return nil
}

// 生成运单
func (this *LogisticsService) AddOrder(ctx context.Context, opt CreateExpressOrderOption) (response *weapp.CreateExpressOrderResponse, err error) {

	result := <-this.orderRep.FindById(ctx, opt.OrderId)
	if result.Error != nil {
		return
	}
	order := result.Result.(*models.Order)

	deliveries, err := this.GetAllDelivery()
	if err != nil {
		return nil, err
	}
	var delivery *models.LogisticsInfo
	var deliveryService *models.LogisticsInfoService
	for _, info := range deliveries {
		if info.DeliveryId == opt.DeliveryId {
			delivery = info
			break
		}
	}
	for _, s := range delivery.Services {
		if s.Type == opt.DeliveryServiceType {
			deliveryService = &s
			break
		}
	}

	if delivery == nil {
		// 物流公司未注册到微信小程序物流助手
		return nil, err2.Err422.F("DeliveryId[%s],不合法", opt.DeliveryId)
	}
	if deliveryService == nil {
		// 物流公司服务不匹配
		return nil, err2.Err422.F("DeliveryServiceType[%s],不合法", opt.DeliveryServiceType)
	}
	if delivery.BizID == "" {
		// 缺少物流结账号
		return nil, err2.Err422.F("%s缺少 biz_id", delivery.DeliveryName)
	}

	expressOrder := weapp.ExpressOrder{
		OrderID:      order.OrderNo,           // 订单ID，须保证全局唯一，不超过512字节
		OpenID:       order.User.WechatMiniId, // 用户openid，当add_source=2时无需填写（不发送物流服务通知）
		DeliveryID:   opt.DeliveryId,          // 快递公司ID，参见getAllDelivery
		BizID:        delivery.BizID,          // 快递客户编码或者现付编码
		CustomRemark: opt.CustomRemark,        // 快递备注信息，比如"易碎物品"，不超过1024字节

		Receiver: weapp.ExpreseeUserInfo{
			Name:     order.UserAddress.ContactName,
			Mobile:   order.UserAddress.ContactPhone,
			Province: order.UserAddress.Province,
			City:     order.UserAddress.City,
			Area:     order.UserAddress.Areas,
			Address:  order.UserAddress.Addr,
		}, // 收件人信息

		Shop: weapp.ExpressShop{
			WXAPath:    fmt.Sprintf("pages/home/order/detail?id=%s", order.GetID()),
			IMGUrl:     order.GetAvatar(),
			GoodsName:  order.GoodsName(30),
			GoodsCount: uint(order.ItemCount),
		},
		Insured: weapp.ExpressInsure{
			Used:  weapp.InsureStatus(opt.UseInsured),
			Value: uint(opt.InsuredValue),
		},
		Service: weapp.ExpressService{
			Type: deliveryService.Type,
			Name: deliveryService.Name,
		},
	}

	// 发件人信息
	sender := weapp.ExpreseeUserInfo{
		Name:     opt.Shop.Address.Name,
		Province: opt.Shop.Address.Province,
		City:     opt.Shop.Address.City,
		Area:     opt.Shop.Address.Areas,
		Address:  opt.Shop.Address.Addr,
	}
	if opt.Shop.Address.Tel != "" {
		sender.Tel = opt.Shop.Address.Tel
	} else {
		sender.Mobile = opt.Shop.Address.Phone
	}

	expressOrder.Sender = sender

	cargo := weapp.ExpressCargo{
		Count:  uint(len(opt.Items)),
		Weight: opt.Weight,
		SpaceX: opt.SpaceX,
		SpaceY: opt.SpaceY,
		SpaceZ: opt.SpaceZ,
	}

	for _, item := range opt.Items {
		findItem := order.FindItem(item.ItemId)
		cargo.DetailList = append(cargo.DetailList, weapp.CargoDetail{
			Name:  findItem.Item.Code,
			Count: uint(item.Count),
		})
	}

	expressOrder.Cargo = cargo

	addExpressOrder, err := wechat.SDK.AddExpressOrder(wechat.CreateExpressOrderOption{
		ExpressOrder: &expressOrder,
		AddSource:    0,
		ExpectTime:   uint(opt.ExpectTime.Unix()),
	}, true)

	if err != nil {
		return nil, err
	}

	var waybillData []*models.WaybillDataItem

	for _, billData := range addExpressOrder.WaybillData {
		waybillData = append(waybillData, &models.WaybillDataItem{
			Key:   billData.Key,
			Value: billData.Value,
		})
	}

	var items []*models.LogisticsItem

	for _, o := range opt.Items {
		items = append(items, &models.LogisticsItem{
			ItemId: o.ItemId,
			Count:  o.Count,
			ShopId: opt.Shop.GetID(),
		})
	}

	l := &models.LogisticsOption{
		NoDelivery:  false,
		DeliveryId:  opt.DeliveryId,
		ShopId:      opt.Shop.GetID(),
		Items:       items,
		WaybillID:   addExpressOrder.WaybillID,
		WaybillData: waybillData,
	}

	// 写入订单物流
	if err := order.Shipment(l); err != nil {
		// 写入订单物流异常，
		// 取消发货单
		_, err := this.CancelOrder(wechat.CancelExpressOrderOption{
			OrderId:    order.OrderNo,
			OpenId:     order.User.WechatMiniId,
			DeliveryId: opt.DeliveryId,
			WaybillId:  addExpressOrder.WaybillID,
		})
		if err != nil {
			return nil, err
		}

	}

	// 保存订单
	saved := <-this.orderRep.Save(ctx, order)
	if saved.Error != nil {
		// 保存订单失败
		// 取消发货运单
		_, err := this.CancelOrder(wechat.CancelExpressOrderOption{
			OrderId:    order.OrderNo,
			OpenId:     order.User.WechatMiniId,
			DeliveryId: opt.DeliveryId,
			WaybillId:  addExpressOrder.WaybillID,
		})
		if err != nil {
			return nil, err
		}
	}
	return addExpressOrder, nil
}

// 取消运单
func (this *LogisticsService) CancelOrder(opt wechat.CancelExpressOrderOption) (*weapp.CommonError, error) {
	order, err := wechat.SDK.CancelOrder(opt)
	if err != nil {
		return nil, err
	}
	if err := order.GetResponseError(); err != nil {
		return nil, err
	}
	return order, nil
}
