package services

import (
	"context"
	"errors"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/wechat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type RefundService struct {
	orderResp *repositories.OrderRep
	refundRep *repositories.RefundRep
}

func NewRefundService(refundRep *repositories.RefundRep, orderResp *repositories.OrderRep) *RefundService {
	return &RefundService{refundRep: refundRep, orderResp: orderResp}
}

type RefundOption struct {
	RefundId string `json:"refund_id" form:"refund_id"`
	Desc     string `json:"desc"`
	Close    bool   `json:"close"`
}

func (refundOpt *RefundOption) IsValid() error {
	if refundOpt.RefundId == "" {
		return errors.New("RefundId is empty")
	}
	return nil
}

type MakeRefundOption struct {
	Order *models.Order           `json:"-" form:"-"`
	Items []*MakeRefundItemOption `json:"items"`
	Type  uint                    `json:"type"`
	Desc  string                  `json:"desc"` // 退款原因
}

type MakeRefundItemOption struct {
	ItemId       string `json:"item_id" form:"item_id"`             // 对应 itemId
	RefundAmount uint64 `json:"refund_amount" form:"refund_amount"` // 申请退款总额
	Count        uint64 `json:"count" form:"count"`                 // 申请退款数量
}

func (opt *MakeRefundOption) isValid(refunds []*models.Refund) error {
	rs := models.Refunds(refunds)
	for _, item := range opt.Items {
		orderItem := opt.Order.FindItem(item.ItemId)
		if orderItem == nil {
			return err2.Err422.F("退款产品与订单不匹配")
		}
		qty, amount := rs.CountItemQtyAndAmount(item.ItemId)
		if item.Count > uint64(orderItem.Count)-qty {
			return err2.Err422.F("退款产品数量[%d]溢出，退款中数量(包含已退款)[%d]，实际购买数量[%d]", item.Count, qty, orderItem.Count)
		}
		if item.RefundAmount > uint64(orderItem.TotalAmount)-amount {
			return err2.Err422.F("item[%s]退款金额[%s]溢出，退款中金额(包含已退款)[%s]，该商品实际支付金额[%s]", item.ItemId, utils.ToMoneyString(item.RefundAmount), utils.ToMoneyString(amount), utils.ToMoneyString(orderItem.TotalAmount))
		}
	}
	return nil
}

// 列表
func (srv *RefundService) Pagination(ctx context.Context, req *request.IndexRequest) (refunds []*models.Refund, pagination response.Pagination, err error) {
	results := <-srv.refundRep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Refund), results.Pagination, nil
}

// 获取退款单
func (srv *RefundService) FindById(ctx context.Context, id string) (refund *models.Refund, err error) {
	result := <-srv.refundRep.FindById(ctx, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return result.Result.(*models.Refund), nil
}

// 获取订单关联的退款单
func (srv *RefundService) FindRefundByOrderId(ctx context.Context, orderId string) (refunds []*models.Refund, err error) {
	return srv.refundRep.FindMany(ctx, bson.M{
		"order_id": orderId,
	})
}

// 申请退款
func (srv *RefundService) MakeRefund(ctx context.Context, opt *MakeRefundOption) (refund *models.Refund, err error) {
	if opt.Order.Status == models.OrderStatusPrePay || opt.Order.Status == models.OrderStatusFailed {
		return nil, err2.Err422.F("当前订单状态处于[%s]，不允许进行退款操作！", opt.Order.StatusText())
	}
	// 查询当前订单，对应的所有的退款单(除关闭的)
	refunds, err := srv.refundRep.FindMany(ctx, bson.M{"order_on": opt.Order.OrderNo, "status": bson.M{"$nin": models.RefundFailedStatus}})
	if err != nil {
		return nil, err
	}
	if err := opt.isValid(refunds); err != nil {
		return nil, err
	}
	order := opt.Order
	// 生成退款单
	refundItems := make([]*models.RefundItem, 0)
	var totalAmount uint64
	for _, item := range opt.Items {
		orderItem := order.FindItem(item.ItemId)
		orderItem.Refunding = true
		refundItem := &models.RefundItem{
			ItemId:      item.ItemId,
			Qty:         item.Count,
			TotalAmount: item.RefundAmount,
			Item:        orderItem.Item,
		}
		totalAmount += item.RefundAmount
		refundItems = append(refundItems, refundItem)
	}
	refund = &models.Refund{
		RefundNo:    utils.RandomRefundOrderNo(""),
		OrderId:     order.GetID(),
		OrderNo:     order.OrderNo,
		PaymentNo:   order.Payment.PaymentNo,
		OpenId:      order.User.WechatMiniId,
		ReturnCode:  "",
		Status:      models.RefundStatusApply,
		TotalAmount: totalAmount,
		Items:       refundItems,
		RefundDesc:  opt.Desc,
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
		created := <-srv.refundRep.Create(sessionContext, refund)
		if created.Error != nil {
			session.AbortTransaction(sessionContext)
			log.Errorf("创建退款单失败,err:%s", created.Error)
			return created.Error
		}

		refund = created.Result.(*models.Refund)

		saved := <-srv.orderResp.Save(sessionContext, order)

		if saved.Error != nil {
			session.AbortTransaction(sessionContext)
			log.Errorf("更新订单退款状态失败,err:%s", saved.Error)
			return created.Error
		}

		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	return refund, err
}

// 同意退款/进行退款(目前做的整单退，如果要做单个的话需要调整RefundOption添加item)
func (srv *RefundService) AgreeRefund(ctx context.Context, refundId string) (refund *models.Refund, err error) {
	result := <-srv.refundRep.FindById(ctx, refundId)
	if result.Error != nil {

	}
	refund = result.Result.(*models.Refund)
	if refund.Status != models.RefundStatusApply {
		err2.Err422.F("refund no[%s]，不合法,当前退款单状态[%d],需要状态为[%d]", refund.RefundNo, refund.Status, models.RefundStatusApply)
	}

	order, err := srv.orderResp.FindByOrderNo(ctx, refund.OrderNo)
	if err != nil {
		return
	}
	if order.RefundChannel == false {
		// 退款通道已关闭
	}

	refundOpt := &wechat.RefundOption{
		OutTradeNo:  refund.OrderNo,
		OutRefundNo: refund.RefundNo,
		TotalFee:    order.ActualAmount,
		RefundFee:   refund.TotalAmount,
		RefundDesc:  refund.RefundDesc,
	}

	response, err := wechat.Pay.Refund(refundOpt)
	spew.Dump(response)
	if err != nil {
		// 参数错误
		return nil, err
	}

	if response.ReturnCode == "SUCCESS" {
		// 退款申请审核通过
		refund.Status = models.RefundStatusAgreed
		// 接口调用成功
		if response.ResultCode == "SUCCESS" {
			// 退款申请接收成功
			refund.Status = models.RefundStatusRefunding
			// 删除失败记录
			refund.FailedLog = nil

		} else {
			// 提交业务失败
			log.Errorf("refund_no[%s]，提交退款失败！，错误码[%s]，失败原因[%s]", refund.RefundNo, response.ErrCode, response.ErrCodeDes)
			// 记录失败原因
			refund.FailedLog = &models.FailedLog{
				ErrCode:    response.ErrCode,
				ErrCodeDes: response.ErrCodeDes,
			}
			err = err2.Err422.F("提交退款失败！，错误码[%s]，失败原因[%s]", response.ErrCode, response.ErrCodeDes)
		}
		saved := <-srv.refundRep.Save(ctx, refund)
		if saved.Error != nil {
			log.Errorf("agree refund,save refund error:%s", saved.Error)
			return nil, saved.Error
		}
		refund = saved.Result.(*models.Refund)
	}

	return refund, err
}

// 拒绝退款
func (srv *RefundService) RejectRefund(ctx context.Context, authenticatable auth.Authenticatable, opt *RefundOption) (refund *models.Refund, err error) {
	if err := opt.IsValid(); err != nil {
		return nil, err
	}
	results := <-srv.refundRep.FindById(ctx, opt.RefundId)
	if results.Error != nil {
		return nil, results.Error
	}
	refund = results.Result.(*models.Refund)

	refund.Status = models.RefundStatusClosed
	refund.RejectDesc = opt.Desc

	var canceler *models.RefundCanceler

	if admin, ok := authenticatable.(*models.Admin); ok {
		canceler = &models.RefundCanceler{
			Type:   "admin",
			Id:     admin.GetID(),
			Name:   admin.Nickname,
			Avatar: "",
		}
	}
	if user, ok := authenticatable.(*models.User); ok {
		canceler = &models.RefundCanceler{
			Type:   "user",
			Id:     user.GetID(),
			Name:   user.Nickname,
			Avatar: user.Avatar,
		}
	}
	refund.Canceler = canceler

	queryResults := <-srv.orderResp.FindById(ctx, refund.OrderId)
	if queryResults.Error != nil {
		err = queryResults.Error
		return
	}
	order := queryResults.Result.(*models.Order)
	order.RefundChannel = opt.Close
	refunds, _ := srv.refundRep.FindMany(ctx, bson.M{"order_no": refund.OrderNo, "status": bson.M{"$in": models.RefundingStatus}})
	rs := make(models.Refunds, 0)
	for _, r := range refunds {
		if r.GetID() != refund.GetID() {
			rs = append(rs, r)
		}
	}

	for _, item := range refund.Items {
		findItem := order.FindItem(item.ItemId)
		if findItem != nil {
			var qty uint64
			var amount uint64
			for _, r := range rs {
				refundItem := r.FindItem(item.ItemId)
				if refundItem != nil {
					qty += refundItem.Qty
					amount += refundItem.TotalAmount
				}
			}
			if qty > 0 || amount > 0 {
				findItem.Refunding = true
			} else {
				findItem.Refunding = false
			}
		}
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
		saved := <-srv.refundRep.Save(sessionContext, refund)
		if saved.Error != nil {
			// 保存退款单失败
			session.AbortTransaction(sessionContext)
			return saved.Error
		}
		refund = saved.Result.(*models.Refund)

		orderSaved := <-srv.orderResp.Save(sessionContext, order)
		if orderSaved.Error != nil {
			// 保存失败
			session.AbortTransaction(sessionContext)
			return orderSaved.Error
		}
		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	return refund, err
}

// 退款回调
func (srv *RefundService) RefundNotify(ctx context.Context, req *http.Request) (refund *models.Refund, err error) {
	// parse
	notifyReq, err := wechat.Pay.ParseRefundNotifyResult(req)
	if err != nil {
		return nil, err
	}
	spew.Dump(notifyReq)
	// deal with
	refundNo := notifyReq.OutRefundNo

	refund, err = srv.refundRep.FindByRefundNo(ctx, refundNo)
	if err != nil {
		// 退款单不存在
		return nil, err
	}

	if refund.ReturnCode == "SUCCESS" {
		// 已经标记为退款成功
		// 不做处理
		return nil, nil
	}

	if refund.ReturnCode == "REFUNDCLOSE" {
		// 当前退款单已经处于退款关闭状态，不做处理
		return nil, nil
	}

	order, err := srv.orderResp.FindByOrderNo(ctx, refund.OrderNo)
	if err != nil {
		return nil, err
	}

	refunds, _ := srv.refundRep.FindMany(ctx, bson.M{"order_no": refund.OrderNo, "status": bson.M{"$in": models.RefundingStatus}})
	rs := make(models.Refunds, 0)
	for _, r := range refunds {
		if r.GetID() != refund.GetID() {
			rs = append(rs, r)
		}
	}

	refund.ReturnCode = notifyReq.RefundStatus
	if notifyReq.RefundStatus == "SUCCESS" && refund.Status == models.RefundStatusRefunding {
		refund.Status = models.RefundStatusDone
		var qty uint64
		var amount uint64
		for _, item := range refund.Items {
			findItem := order.FindItem(item.ItemId)
			if findItem != nil {
				findItem.RemainderAmount -= item.TotalAmount
				findItem.RemainderQty -= item.Qty
				for _, r := range rs {
					refundItem := r.FindItem(item.ItemId)
					if refundItem != nil {
						qty += refundItem.Qty
						amount += refundItem.TotalAmount
					}
				}
				if qty > 0 || amount > 0 {
					findItem.Refunding = true
				} else {
					findItem.Refunding = false
				}
				qty += item.Qty
				amount += item.TotalAmount
			}
		}
		order.RemainderAmount -= refund.TotalAmount
		if qty == order.ItemCount && amount == order.ActualAmount {
			order.Status = models.OrderStatusFailed
			order.RefundMark = 2
		} else {
			order.RefundMark = 1
		}
	}
	if notifyReq.RefundStatus == "REFUNDCLOSE" && refund.Status == models.RefundStatusRefunding {
		refund.Status = models.RefundStatusClosed
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
		// 保存退款信息
		saved := <-srv.refundRep.Save(ctx, refund)
		if saved.Error != nil {
			// 保存失败
			session.AbortTransaction(sessionContext)
			return saved.Error
		}
		refund = saved.Result.(*models.Refund)
		spew.Dump(order)
		orderSaved := <-srv.orderResp.Save(sessionContext, order)
		if orderSaved.Error != nil {
			// 保存失败
			session.AbortTransaction(sessionContext)
			return orderSaved.Error
		}
		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	return refund, err
}

// 取消退款
func (srv *RefundService) CancelRefund(ctx context.Context, opt *RefundOption, authenticatable auth.Authenticatable, closeChannel bool) (refund *models.Refund, err error) {
	// 用户取消退款，取消后关闭退款通道
	if err := opt.IsValid(); err != nil {
		return nil, err
	}

	results := <-srv.refundRep.FindById(ctx, opt.RefundId)
	if results.Error != nil {
		return nil, results.Error
	}

	refund = results.Result.(*models.Refund)
	if !refund.CanCancel() {
		// 当前状态不允许关闭退款
		return nil, err2.Err422.F("当前状态无法关闭退款单")
	}

	refund.Status = models.RefundStatusClosed
	canceler := &models.RefundCanceler{}
	if user, ok := authenticatable.(*models.User); ok {
		canceler.Type = "user"
		canceler.Id = user.GetID()
		canceler.Name = user.Nickname
		canceler.Avatar = user.Avatar
		closeChannel = true
	}
	if admin, ok := authenticatable.(*models.Admin); ok {
		canceler.Type = "admin"
		canceler.Id = admin.GetID()
		canceler.Name = admin.Nickname
	}
	// 记录关闭退款单操作者
	refund.Canceler = canceler
	// 删除失败退款记录
	refund.FailedLog = nil

	queryResults := <-srv.orderResp.FindById(ctx, refund.OrderId)
	if queryResults.Error != nil {
		err = queryResults.Error
		return
	}
	order := queryResults.Result.(*models.Order)

	refunds, _ := srv.refundRep.FindMany(ctx, bson.M{"order_no": refund.OrderNo, "status": bson.M{"$in": models.RefundingStatus}})
	rs := make(models.Refunds, 0)
	for _, r := range refunds {
		if r.GetID() != refund.GetID() {
			rs = append(rs, r)
		}
	}

	for _, item := range refund.Items {
		findItem := order.FindItem(item.ItemId)
		if findItem != nil {
			var qty uint64
			var amount uint64
			for _, r := range rs {
				refundItem := r.FindItem(item.ItemId)
				if refundItem != nil {
					qty += refundItem.Qty
					amount += refundItem.TotalAmount
				}
			}
			if qty > 0 || amount > 0 {
				findItem.Refunding = true
			} else {
				findItem.Refunding = false
			}
		}
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
		saved := <-srv.refundRep.Save(sessionContext, refund)
		if saved.Error != nil {
			// 保存退款单失败
			session.AbortTransaction(sessionContext)
			return saved.Error
		}
		refund = saved.Result.(*models.Refund)

		orderSaved := <-srv.orderResp.Save(sessionContext, order)
		if orderSaved.Error != nil {
			// 保存失败
			session.AbortTransaction(sessionContext)
			return orderSaved.Error
		}
		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	return refund, err
}

// 取消退款
func (srv *RefundService) Cancel(ctx context.Context, entity *models.Refund, authenticatable auth.Authenticatable) (refund *models.Refund, err error) {
	if !entity.CanCancel() {
		// 当前状态不允许关闭退款
		return nil, err2.Err422.F("当前状态无法关闭退款单")
	}

	entity.Status = models.RefundStatusClosed
	canceler := &models.RefundCanceler{}
	if user, ok := authenticatable.(*models.User); ok {
		canceler.Type = "user"
		canceler.Id = user.GetID()
		canceler.Name = user.Nickname
		canceler.Avatar = user.Avatar
	}
	if admin, ok := authenticatable.(*models.Admin); ok {
		canceler.Type = "admin"
		canceler.Id = admin.GetID()
		canceler.Name = admin.Nickname
	}
	// 记录关闭退款单操作者
	entity.Canceler = canceler

	queryResults := <-srv.orderResp.FindById(ctx, entity.OrderId)
	if queryResults.Error != nil {
		err = queryResults.Error
		return
	}
	order := queryResults.Result.(*models.Order)

	refunds, _ := srv.refundRep.FindMany(ctx, bson.M{"order_no": entity.OrderNo, "status": bson.M{"$in": models.RefundingStatus}})
	rs := make(models.Refunds, 0)
	for _, r := range refunds {
		if r.GetID() != entity.GetID() {
			rs = append(rs, r)
		}
	}

	for _, item := range entity.Items {
		findItem := order.FindItem(item.ItemId)
		if findItem != nil {
			var qty uint64
			var amount uint64
			for _, r := range rs {
				refundItem := r.FindItem(item.ItemId)
				if refundItem != nil {
					qty += refundItem.Qty
					amount += refundItem.TotalAmount
				}
			}
			if qty > 0 || amount > 0 {
				findItem.Refunding = true
			} else {
				findItem.Refunding = false
			}
		}
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
		saved := <-srv.refundRep.Save(sessionContext, entity)
		if saved.Error != nil {
			// 保存退款单失败
			session.AbortTransaction(sessionContext)
			return saved.Error
		}
		refund = saved.Result.(*models.Refund)

		orderSaved := <-srv.orderResp.Save(sessionContext, order)
		if orderSaved.Error != nil {
			// 保存失败
			session.AbortTransaction(sessionContext)
			return orderSaved.Error
		}
		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	return refund, err
}
