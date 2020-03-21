package services

import (
	"context"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// create struct
type OrderCreateOption struct {
	UserAddress  orderUserAddress         `json:"user_address" form:"user_address"`
	TakeGoodType int                      `json:"take_good_type" form:"take_good_type"`
	OrderItems   []*OrderItemCreateOption `json:"order_items" form:"order_items"`
	OrderAmount  uint64                   `json:"order_amount" form:"order_amount"`
	ActualAmount uint64                   `json:"actual_amount" form:"actual_amount"`
}

type orderUserAddress struct {
	Id           string `json:"id" form:"id"`
	ContactName  string `json:"contact_name" form:"contact_name"`
	ContactPhone string `json:"contact_phone" form:"contact_phone"`
	Province     string `json:"province" form:"province"`
	City         string `json:"city" form:"city"`
	Areas        string `json:"areas" form:"areas"`
	Addr         string `json:"addr" form:"addr"`
}

type OrderItemCreateOption struct {
	ItemId            string `json:"item_id" form:"item_id"`
	ProductId         string
	Qty               int64    `json:"qty" form:"qty"`                             // 购买数量
	Price             int64    `json:"price" form:"price"`                         // 商品价格
	MutexPromotion    *string  `json:"mutexPromotion" form:"mutexPromotion"`       // 参加的互斥活动，互斥活动只允许同时参加一个
	UnMutexPromotions []string `json:"unMutexPromotions" form:"unMutexPromotions"` // 参加的非互斥活动
}

func (opt *OrderCreateOption) IsValid() error {
	// user address
	if opt.UserAddress.ContactName == "" {
		return err2.Err422.F("empty contact name")
	}
	if opt.UserAddress.ContactPhone == "" {
		return err2.Err422.F("empty contact phone")
	}
	if opt.UserAddress.Province == "" {
		return err2.Err422.F("empty province")
	}
	if opt.UserAddress.City == "" {
		return err2.Err422.F("empty city")
	}
	if opt.UserAddress.Areas == "" {
		return err2.Err422.F("empty areas")
	}
	if opt.UserAddress.Addr == "" {
		return err2.Err422.F("empty address")
	}
	// items
	if len(opt.OrderItems) == 0 {
		return err2.Err422.F("empty order items")
	}
	for _, item := range opt.OrderItems {
		if item.ItemId == "" {
			return err2.Err422.F("empty item id")
		}
		if item.Qty == 0 {
			return err2.Err422.F("invalid item count")
		}
	}
	// amount
	if opt.OrderAmount == 0 {
		return err2.Err422.F("invalid order amount")
	}
	return nil
}

// deliver struct
type DeliverOption struct {
	Options []*models.LogisticsOption `json:"options" form:"options"`
}

// confirm struct
type ConfirmOption struct {
	OrderNo string `json:"order_no" form:"order_no"`
}

func (opt *ConfirmOption) IsValid() error {
	if opt.OrderNo == "" {
		return err2.Err422.F("empty order no")
	}
	return nil
}

// cancel struct
type CancelOption struct {
	OrderNo string
}

func (opt *CancelOption) IsValid() error {
	if opt.OrderNo == "" {
		return err2.Err422.F("empty order no")
	}
	return nil
}

type OrderService struct {
	orderRep   *repositories.OrderRep
	commentRep *repositories.CommentRep
	//orderInventoryLogRep *repositories.OrderInventoryLogRep
	promotionSrv *PromotionService
	productSrv   *ProductService
}

func NewOrderService(orderRep *repositories.OrderRep, promotionSrv *PromotionService, productSrv *ProductService) *OrderService {
	return &OrderService{orderRep: orderRep, promotionSrv: promotionSrv, productSrv: productSrv}
}

// 列表
func (srv *OrderService) Pagination(ctx context.Context, req *request.IndexRequest) (orders []*models.Order, pagination response.Pagination, err error) {
	results := <-srv.orderRep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Order), results.Pagination, nil
}

// 详情
func (srv *OrderService) FindById(ctx context.Context, id string) (order *models.Order, err error) {
	results := <-srv.orderRep.FindById(ctx, id)
	if results.Error != nil {
		return nil, results.Error
	}
	return results.Result.(*models.Order), nil
}

// 创建订单
func (srv *OrderService) Create(ctx context.Context, userInfo *models.User, opt *OrderCreateOption) (order *models.Order, err error) {
	// 校验数据: 数据有效 -> 产品有效 -> 库存充足 -> 金额匹配
	// 数据有效
	if err = opt.IsValid(); err != nil {
		return nil, err
	}
	// 校验产品有效且价格合法且库存充足
	var calcAmount uint64 = 0
	detailItem := make([]*models.OrderItem, 0)
	for _, orderItem := range opt.OrderItems {
		// todo: here should be optimized
		// 价格合法
		item, err := srv.productSrv.FindItemById(ctx, orderItem.ItemId)
		if err != nil {
			return nil, err
		}
		orderItem.ProductId = item.Product.Id
		// valid price
		if item.PromotionPrice != orderItem.Price {
			return nil, err2.Err422.F(fmt.Sprintf("invalid item price with %s-%d", orderItem.ItemId, orderItem.Price))
		}
		// 库存充足
		if item.Qty < orderItem.Qty {
			return nil, err2.Err422.F(fmt.Sprintf("item %s inventory not enough which remain %d", orderItem.ItemId, item.Qty))
		}
		// calculate all amount
		calcAmount += uint64(item.PromotionPrice * orderItem.Qty)
		// store item with detail
		detailItem = append(detailItem, &models.OrderItem{
			Item:  item.ToAssociated(),
			Count: orderItem.Qty,
			Price: orderItem.Price,
		})
	}
	// 金额匹配
	if calcAmount != opt.OrderAmount {
		return nil, err2.Err422.F("order amount not equal")
	}
	// 计算优惠
	promotionResult := srv.getDiscounts(ctx, opt)

	if opt.OrderAmount-uint64(promotionResult.SalePrices) != opt.ActualAmount {
		return nil, err2.Err422.F("invalid order actual amount")
	}
	// 生成订单
	order = srv.generateOrder(userInfo, opt, detailItem, promotionResult)
	// save order into db
	// transaction of mongo
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return nil, err
	}
	if err = session.StartTransaction(); err != nil {
		return nil, err
	}
	// create order and deduct/lock inventory
	var orderRes *models.Order
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// define
		// itemInventoryLog := make([]*models.ItemInventoryLog, 0)
		// create order
		created := <-srv.orderRep.Create(sessionContext, order)
		if created.Error != nil {
			session.AbortTransaction(sessionContext)
			return created.Error
		}
		orderRes = created.Result.(*models.Order)
		// deduct / lock inventory
		//addressLocation, _ := order.UserAddress.Location()
		for _, orderItem := range order.OrderItems {
			// 直接扣库存
			if err := srv.productSrv.ItemService.DecQty(sessionContext, orderItem.Item.Id, orderItem.Count); err != nil {
				// 扣库存失败
				session.AbortTransaction(sessionContext)
				return err
			}
			//itemInventory, lockErr := srv.inventoryRep.LockByItemId(sessionContext, orderItem.Item.Id, orderItem.Count, 0, addressLocation)
			//if lockErr != nil {
			//	session.AbortTransaction(sessionContext)
			//	return lockErr
			//}
			// fill item inventory log
			//itemInventoryLog = append(itemInventoryLog, &models.ItemInventoryLog{
			//	ItemId: orderItem.Item.Id,
			//	//InventoryId: itemInventory.GetID(),
			//})
		}
		// create order inventory log
		//orderInventoryLog := &models.OrderInventoryLog{
		//	OrderNo:          order.OrderNo,
		//	ItemInventoryLog: itemInventoryLog,
		//}
		//logCreated := <-srv.orderInventoryLogRep.Create(sessionContext, orderInventoryLog)
		//if logCreated.Error != nil {
		//	session.AbortTransaction(sessionContext)
		//	return logCreated.Error
		//}
		// return
		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	// return
	if err != nil {
		return nil, err
	}
	return orderRes, nil
}

// 计算促销优惠
func (srv *OrderService) getDiscounts(ctx context.Context, opt *OrderCreateOption) *models.PromotionResult {
	return srv.promotionSrv.CalculateByOrder(ctx, opt.OrderItems)
}

// 生成订单
func (srv *OrderService) generateOrder(user *models.User, opt *OrderCreateOption, orderItems []*models.OrderItem, promotionResult *models.PromotionResult) *models.Order {
	if detail := promotionResult.Detail(); detail != nil {
		for _, item := range orderItems {
			if info := detail.FindByItemId(item.Item.Id); info != nil {
				item.PromotionInfo = info
				item.Amount = item.Price - info.UnitSalePrices
			} else {
				item.Amount = item.Price
			}
		}
	} else {
		for _, item := range orderItems {
			item.Amount = item.Price
		}
	}

	var itemCount int64
	for _, orderItem := range opt.OrderItems {
		itemCount += orderItem.Qty
	}
	// build model of order
	resOrder := &models.Order{
		OrderNo:      utils.RandomOrderNo(""),
		ItemCount:    itemCount,
		OrderAmount:  opt.OrderAmount,
		ActualAmount: opt.ActualAmount,
		OrderItems:   orderItems,
		User:         user.ToAssociated(),
		UserAddress: &models.AssociatedUserAddress{
			Id:           opt.UserAddress.Id,
			ContactName:  opt.UserAddress.ContactName,
			ContactPhone: opt.UserAddress.ContactPhone,
			Province:     opt.UserAddress.Province,
			City:         opt.UserAddress.City,
			Areas:        opt.UserAddress.Areas,
			Addr:         opt.UserAddress.Addr,
		},
		TakeGoodType:  opt.TakeGoodType,
		Logistics:     nil,                        // todo: confirm how different between using nil and &models.Logistics
		Payment:       nil,                        // todo: same with above
		PromotionInfo: promotionResult.Overview(), // 促销总览
		Status:        models.OrderStatusPrePay,
	}
	// return
	return resOrder
}

// 发货
func (srv *OrderService) Deliver(ctx context.Context, order *models.Order, opt *DeliverOption) (model *models.Order, err error) {
	if err := order.Shipment(opt.Options...); err != nil {
		return nil, err
	}
	saved := <-srv.orderRep.Save(ctx, order)
	if saved.Error != nil {
		return nil, saved.Error
	}
	return saved.Result.(*models.Order), nil
}

// 确认收货
func (srv *OrderService) Confirm(ctx context.Context, opt *ConfirmOption) error {
	if err := opt.IsValid(); err != nil {
		return err
	}
	// 获取用户
	authUser := ctx2.GetUser(ctx)
	userInfo, ok := authUser.(*models.User)
	if !ok {
		return err2.Err422.F("invalid user who is unauthenticated")
	}
	// 查询订单并校验状态
	orderRes := <-srv.orderRep.FindOne(ctx, map[string]interface{}{
		"order_no": opt.OrderNo,
	})
	if orderRes.Error != nil {
		return orderRes.Error
	}
	order := orderRes.Result.(*models.Order)
	if order.Status != models.OrderStatusPreConfirm {
		return err2.Err422.F(fmt.Sprintf("order %s can not be comfirm caused of not pre confirm status", opt.OrderNo))
	}
	// 校验是否为用户本人
	if order.User.Id != userInfo.GetID() {
		return err2.Err422.F(fmt.Sprintf("order %s can not be comfirm caused of invalid user", opt.OrderNo))
	}
	// 更新
	updated := <-srv.orderRep.Update(ctx, order.GetID(), bson.M{
		"$set": bson.M{
			"status": models.OrderStatusPreEvaluate,
		},
	})
	if updated.Error != nil {
		return updated.Error
	}
	// return
	return nil
}

type OrderCommentOption struct {
	Rate      float64        `json:"rate"`
	ProductId string         `json:"product_id" form:"product_id"`
	ItemId    string         `json:"item_id" form:"item_id"`
	Content   string         `json:"content"`
	Images    []*qiniu.Image `json:"images"`
}

// 评价
func (srv *OrderService) Comment(ctx context.Context, order *models.Order, user *models.User, opts []*OrderCommentOption) error {
	// 验证
	if order.User.Id != user.GetID() {
		return err2.Err422.F("评论失败")
	}
	// 状态
	if !order.CanComment() {
		return err2.Err422.F("评论失败")
	}
	var comments []*models.Comment
	for _, opt := range opts {
		if find := order.FindItem(opt.ItemId); find == nil {
			return err2.Err422.F("评论失败")
		}
		comments = append(comments, &models.Comment{
			ProductId: opt.ProductId,
			ItemId:    opt.ItemId,
			OrderId:   order.GetID(),
			User: &models.CommentUser{
				UserId:   user.GetID(),
				Avatar:   user.Avatar.Src(),
				Nickname: user.Nickname,
			},
			Content: opt.Content,
			Images:  nil, // 暂时不存图片
			Rate:    opt.Rate,
		})
	}

	for _, com := range comments {
		created := <-srv.commentRep.Create(ctx, com)
		if created.Error != nil {
			// 写日志
		}
	}
	return nil
}

// 取消订单
func (srv *OrderService) Cancel(ctx context.Context, order *models.Order) error {
	if err := order.StatusToFailed(); err != nil {
		return err
	}
	// 开启事务
	var session mongo.Session
	var err error
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return err
	}
	if err = session.StartTransaction(); err != nil {
		return err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// 退还库存
		for _, item := range order.OrderItems {
			if err := srv.productSrv.ItemService.IncQty(sessionContext, item.Item.Id, item.Count); err != nil {
				session.AbortTransaction(sessionContext)
				return err
			}
		}
		// 保存订单状态
		saved := <-srv.orderRep.Save(sessionContext, order)
		if saved.Error != nil {
			sessionContext.AbortTransaction(sessionContext)
			return saved.Error
		}
		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	return err
}

// 售后工单
func (srv *OrderService) CreateIssue(ctx context.Context, order *models.Order) {

}

// 评论工单
func (srv *OrderService) CommentIssue(ctx context.Context, order *models.Order, user auth.Authenticatable) {

}
