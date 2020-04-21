package services

import (
	"context"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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
	Qty               uint64   `json:"qty" form:"qty"`                             // 购买数量
	Price             uint64   `json:"price" form:"price"`                         // 商品价格
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
	refundRep    *repositories.RefundRep
}

func NewOrderService(orderRep *repositories.OrderRep, promotionSrv *PromotionService, productSrv *ProductService, commentRep *repositories.CommentRep, refundRep *repositories.RefundRep) *OrderService {
	return &OrderService{orderRep: orderRep, promotionSrv: promotionSrv, productSrv: productSrv, commentRep: commentRep, refundRep: refundRep}
}

func (srv *OrderService) Save(ctx context.Context, entity *models.Order) (order *models.Order, err error) {
	saved := <-srv.orderRep.Save(ctx, entity)
	if saved.Error != nil {
		return nil, saved.Error
	}
	return saved.Result.(*models.Order), nil
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

// 通过订单号查询订单
func (srv *OrderService) FindByNo(ctx context.Context, no string) (order *models.Order, err error) {
	return srv.orderRep.FindByOrderNo(ctx, no)
}

// 状态查询
func (srv *OrderService) GetOrderStatus(ctx context.Context, id string) (status int, err error) {
	return srv.orderRep.GetOrderStatus(ctx, id)
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

		// create order
		created := <-srv.orderRep.Create(sessionContext, order)
		if created.Error != nil {
			session.AbortTransaction(sessionContext)
			return created.Error
		}
		orderRes = created.Result.(*models.Order)

		for _, orderItem := range order.OrderItems {
			// 直接扣库存
			if err := srv.productSrv.itemRep.DecQty(sessionContext, orderItem.Item.Id, int64(orderItem.Count)); err != nil {
				// 扣库存失败
				session.AbortTransaction(sessionContext)
				return err
			}
			// 添加销量
			srv.productSrv.UpdateSalesQty(sessionContext, orderItem.Item.Id, int64(orderItem.Count))

		}

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
				payAmount := item.Price*item.Count - info.SalePrices
				// 子订单支付总金额
				item.TotalAmount = payAmount
				item.RemainderAmount = payAmount
			} else {
				payAmount := item.Price * item.Count
				item.TotalAmount = payAmount
				item.RemainderAmount = payAmount
			}
			item.RemainderQty = item.Count
		}
	} else {
		for _, item := range orderItems {
			payAmount := item.Price * item.Count
			item.TotalAmount = payAmount
			item.RemainderAmount = payAmount
			item.RemainderQty = item.Count
		}
	}

	var itemCount uint64
	for _, orderItem := range opt.OrderItems {
		itemCount += orderItem.Qty
	}
	// build model of order
	resOrder := &models.Order{
		OrderNo:         utils.RandomOrderNo(""),
		ItemCount:       itemCount,
		OrderAmount:     opt.OrderAmount,
		ActualAmount:    opt.ActualAmount,
		RemainderAmount: opt.ActualAmount,
		OrderItems:      orderItems,
		User:            user.ToAssociated(),
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
		Logistics:     make([]*models.Logistics, 0), // todo: confirm how different between using nil and &models.Logistics
		Payment:       nil,                          // todo: same with above
		PromotionInfo: promotionResult.Overview(),   // 促销总览
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
func (srv *OrderService) Confirm(ctx context.Context, orderNo string) (order *models.Order, err error) {
	// 查询订单并校验状态
	order, err = srv.orderRep.FindByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	if order.Status != models.OrderStatusPreConfirm {
		return nil, err2.Err422.F("order [%s] can not be comfirm caused of not pre confirm status", orderNo)
	}
	// 更新
	updated := <-srv.orderRep.Update(ctx, order.GetID(), bson.M{
		"$set": bson.M{
			"status":          models.OrderStatusPreEvaluate,
			"refund_channel":  false,
			"comment_channel": true,
		},
	})
	if updated.Error != nil {
		return nil, updated.Error
	}
	// return
	return updated.Result.(*models.Order), nil
}

type OrderCommentOption struct {
	Rate    uint           `json:"rate"`
	Content string         `json:"content"`
	Images  []*qiniu.Image `json:"images"`
}

// 评价
func (srv *OrderService) Comment(ctx context.Context, order *models.Order, user *models.User, opt *OrderCommentOption) (*models.Order, error) {
	var err error
	// 验证
	if order.User.Id != user.GetID() {
		return nil, err2.Err422.F("评论失败")
	}
	// 状态
	if !order.CanComment() {
		return nil, err2.Err422.F("评论失败")
	}
	var comments []*models.Comment

	for _, id := range order.GetProductIds() {
		comments = append(comments, &models.Comment{
			ProductId: id,
			OrderNo:   order.OrderNo,
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

	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return nil, err
	}
	if err = session.StartTransaction(); err != nil {
		return nil, err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err := srv.commentRep.CreateMany(sessionContext, comments); err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}

		// 关闭评论通道,交易完成
		updated := <-srv.orderRep.Update(sessionContext, order.GetID(), bson.M{
			"$set": bson.M{
				"comment_channel": false,
				"commented_at":    time.Now(),
				"status":          models.OrderStatusDone,
			},
		})
		if updated.Error != nil {
			session.AbortTransaction(sessionContext)
			return updated.Error
		}
		order = updated.Result.(*models.Order)
		session.CommitTransaction(sessionContext)
		return nil
	})

	session.EndSession(ctx)

	return order, err
}

// 取消订单
func (srv *OrderService) Cancel(ctx context.Context, order *models.Order, reason string) (model *models.Order, err error) {
	if err := order.StatusToFailed(); err != nil {
		return nil, err
	}
	order.SetCloseReason(reason)
	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return nil, err
	}
	if err = session.StartTransaction(); err != nil {
		return nil, err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// 退还库存
		for _, item := range order.OrderItems {
			if err := srv.productSrv.itemRep.IncQty(sessionContext, item.Item.Id, item.Count); err != nil {
				session.AbortTransaction(sessionContext)
				return err
			}
			// 退回销量
			srv.productSrv.UpdateSalesQty(sessionContext, item.Item.Id, int64(-item.Count))
		}

		// 保存订单状态
		saved := <-srv.orderRep.Save(sessionContext, order)
		if saved.Error != nil {
			sessionContext.AbortTransaction(sessionContext)
			return saved.Error
		}
		order = saved.Result.(*models.Order)
		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	return order, err
}

// 售后工单
func (srv *OrderService) CreateIssue(ctx context.Context, order *models.Order) {

}

// 评论工单
func (srv *OrderService) CommentIssue(ctx context.Context, order *models.Order, user auth.Authenticatable) {

}

// 当天新订单数量
func (srv *OrderService) TodayNewOrderCount(ctx context.Context) int64 {
	result := <-srv.orderRep.Count(ctx,
		bson.M{
			"created_at": bson.M{"$gte": utils.TodayStart(), "$lte": utils.TodayEnd()},
		})
	if result.Error != nil {
		return 0
	}
	return result.Result
}

// 待付款/待发货订单数量
func (srv *OrderService) CountByStatus(ctx context.Context, status int) (count int64, err error) {
	return srv.orderRep.CountByStatus(ctx, status)
}

func (srv *OrderService) AggregateOrderItem(ctx context.Context, req *request.IndexRequest) (res []*models.AggregateOrderItem, pagination response.Pagination, err error) {
	if req.Search != "" {
		req.AppendFilter("order_items", bson.M{"$elemMatch": bson.M{"item.code": primitive.Regex{Pattern: req.Search, Options: "i"}}})
	}
	filters := req.Filters.Unmarshal()
	if shops, ok := filters["shops"]; ok {
		if len(shops.([]interface{})) > 0 {
			req.AppendFilter("logistics.items.shop_id", bson.M{"$in": shops})
		}

	}
	return srv.orderRep.AggregateOrderItem(ctx, req)
}

// 退款中的订单数量
func (srv *OrderService) RefundingCount(ctx context.Context) int64 {
	return srv.orderRep.RefundingCount(ctx)
}
