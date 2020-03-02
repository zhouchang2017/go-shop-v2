package services

import (
	"context"
	"errors"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// create struct
type OrderCreateOption struct {
	UserAddress  orderUserAddress    `json:"user_address" form:"products"`
	OrderItems   []*models.OrderItem `json:"order_items" form:"order_items"`
	OrderAmount  uint64              `json:"order_amount" form:"order_amount"`
	ActualAmount uint64              `json:"actual_amount" form:"actual_amount"`
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

func (opt *OrderCreateOption) IsValid() error {
	// user address
	if opt.UserAddress.ContactName == "" {
		return errors.New("empty contact name")
	}
	if opt.UserAddress.ContactPhone == "" {
		return errors.New("empty contact phone")
	}
	if opt.UserAddress.Province == "" {
		return errors.New("empty province")
	}
	if opt.UserAddress.City == "" {
		return errors.New("empty city")
	}
	if opt.UserAddress.Areas == "" {
		return errors.New("empty areas")
	}
	if opt.UserAddress.Addr == "" {
		return errors.New("empty address")
	}
	// items
	if len(opt.OrderItems) == 0 {
		return errors.New("empty order items")
	}
	// amount
	if opt.OrderAmount == 0 {
		return errors.New("invalid order amount")
	}
	return nil
}

// deliver struct
type DeliverOption struct {
	OrderNo      string `json:"order_no" form:"order_no"`
	OperatorId   string `json:"operator_id" form:"operator_id"`
	OperatorName string `json:"operator_name" form:"operator_name"`
	models.Logistics
}

func (opt *DeliverOption) IsValid() error {
	if opt.OrderNo == "" {
		return errors.New("empty order no")
	}
	if opt.TrackNo == "" {
		return errors.New("empty track no")
	}
	return nil
}

// deliver struct
type ConfirmOption struct {
	OrderNo string `json:"order_no" form:"order_no"`
}

func (opt *ConfirmOption) IsValid() error {
	if opt.OrderNo == "" {
		return errors.New("empty order no")
	}
	return nil
}

type OrderService struct {
	rep *repositories.OrderRep
}

func NewOrderService(rep *repositories.OrderRep) *OrderService {
	return &OrderService{rep: rep}
}

// 列表
func (srv *OrderService) Pagination(ctx context.Context, req *request.IndexRequest) (orders []*models.Order, pagination response.Pagination, err error) {
	results := <-srv.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Order), results.Pagination, nil
}

// 详情
func (srv *OrderService) FindById(ctx context.Context, id string) (order *models.Order, err error) {
	results := <-srv.rep.FindById(ctx, id)
	if results.Error != nil {
		return nil, results.Error
	}
	return results.Result.(*models.Order), nil
}

// 创建订单
func (srv *OrderService) Create(ctx context.Context, opt *OrderCreateOption) (order *models.Order, err error) {
	// 获取用户
	authUser := ctx2.GetUser(ctx)
	userInfo, ok := authUser.(models.User)
	if !ok {
		return nil, errors.New("invalid user who is unauthenticated")
	}
	// 校验数据: 数据有效 - 产品有效 - 库存充足 - 金额匹配
	// 数据有效
	if err = opt.IsValid(); err != nil {
		return nil, err
	}
	// 校验产品有效且库存充足
	var calcAmount uint64 = 0
	// todo
	// 金额匹配
	if calcAmount != opt.OrderAmount {
		return nil, errors.New("order amount not equal")
	}
	// 目前不涉及优惠，暂时取一致
	if discountAmount := srv.getDiscounts(opt); opt.OrderAmount-discountAmount != opt.ActualAmount {
		return nil, errors.New("invalid order actual amount")
	}
	// 生成订单
	order = srv.generateOrder(userInfo, opt)
	// save order into db
	created := <-srv.rep.Create(ctx, &order)
	if created.Error != nil {
		return nil, created.Error
	}
	// return
	return created.Result.(*models.Order), nil
}

func (srv *OrderService) getDiscounts(opt *OrderCreateOption) uint64 {
	return 0
}

func (srv *OrderService) generateOrder(user models.User, opt *OrderCreateOption) *models.Order {
	// build model of order
	resOrder := &models.Order{
		OrderNo:      utils.RandomOrderNo(""),
		ItemCount:    len(opt.OrderItems),
		OrderAmount:  opt.OrderAmount,
		ActualAmount: opt.ActualAmount,
		OrderItems:   opt.OrderItems, // consider to replace with `copy(a, b)`
		User: &models.AssociatedUser{
			Id:       user.GetID(),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Gender:   user.Gender,
		},
		UserAddress: &models.AssociatedUserAddress{
			Id:           opt.UserAddress.Id,
			ContactName:  opt.UserAddress.ContactName,
			ContactPhone: opt.UserAddress.ContactPhone,
			Province:     opt.UserAddress.Province,
			City:         opt.UserAddress.City,
			Areas:        opt.UserAddress.Areas,
			Addr:         opt.UserAddress.Addr,
		},
		Logistics: nil, // todo: confirm how different between using nil and &models.Logistics
		Payment:   nil, // todo: same with above
		Status:    models.OrderStatusPrePay,
	}
	// return
	return resOrder
}

// 发货
func (srv *OrderService) Deliver(ctx context.Context, opt *DeliverOption) error {
	if err := opt.IsValid(); err != nil {
		return err
	}
	// 查询订单并校验状态
	orderRes := <-srv.rep.FindOne(ctx, map[string]interface{}{
		"order_no": opt.OrderNo,
	})
	if orderRes.Error != nil {
		return orderRes.Error
	}
	order := orderRes.Result.(*models.Order)
	if order.Status != models.OrderStatusPreSend {
		return errors.New(fmt.Sprintf("order %s can not be delivered caused of not pre send status", opt.OrderNo))
	}
	// 更新
	updated := <-srv.rep.Update(ctx, order.GetID(), bson.M{
		"$set": bson.M{
			"status": models.OrderStatusPreConfirm,
			"logistics": &models.Logistics{
				Enterprise: opt.Enterprise,
				TrackNo:    opt.TrackNo,
			},
		},
	})
	if updated.Error != nil {
		return updated.Error
	}
	// return
	return nil
}

// 确认收货
func (srv *OrderService) Confirm(ctx context.Context, opt *ConfirmOption) error {
	if err := opt.IsValid(); err != nil {
		return err
	}
	// 获取用户
	authUser := ctx2.GetUser(ctx)
	userInfo, ok := authUser.(models.User)
	if !ok {
		return errors.New("invalid user who is unauthenticated")
	}
	// 查询订单并校验状态
	orderRes := <-srv.rep.FindOne(ctx, map[string]interface{}{
		"order_no": opt.OrderNo,
	})
	if orderRes.Error != nil {
		return orderRes.Error
	}
	order := orderRes.Result.(*models.Order)
	if order.Status != models.OrderStatusPreConfirm {
		return errors.New(fmt.Sprintf("order %s can not be comfirm caused of not pre confirm status", opt.OrderNo))
	}
	// 校验是否为用户本人
	if order.User.Id != userInfo.GetID() {
		return errors.New(fmt.Sprintf("order %s can not be comfirm caused of invalid user", opt.OrderNo))
	}
	// 更新
	updated := <-srv.rep.Update(ctx, order.GetID(), bson.M{
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