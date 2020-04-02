package services

import (
	"context"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/wechat"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type RefundService struct {
	failedRefundRep *repositories.FailedRefundRep
	orderResp       *repositories.OrderRep
}

func NewRefundService(failedRefundRep *repositories.FailedRefundRep, orderResp *repositories.OrderRep) *RefundService {
	return &RefundService{failedRefundRep: failedRefundRep, orderResp: orderResp}
}

type RefundOption struct {
	OrderId  string `json:"order_id" form:"order_id" binding:"required"`
	OrderNo  string `json:"order_no" form:"order_no" binding:"required"`
	RefundNo string `json:"refund_no" form:"refund_no" binding:"required"`
	Desc     string `json:"desc"`
}

func (refundOpt *RefundOption) IsValid() error {
	if refundOpt.OrderId == "" {
		return errors.New("OrderId is empty")
	}
	if refundOpt.OrderNo == "" {
		return errors.New("OrderNo is empty")
	}
	return nil
}

func (srv *RefundService) FindFailedByRefundNo(ctx context.Context, refundNo string) (failedRefund *models.FailedRefund, err error) {
	return srv.failedRefundRep.FindByRefundOn(ctx, refundNo)
}

// 同意退款/进行退款(目前做的整单退，如果要做单个的话需要调整RefundOption添加item)
func (srv *RefundService) AgreeRefund(ctx context.Context, opt *RefundOption) (refundOrder *models.Refund, order *models.Order, err error) {
	if err := opt.IsValid(); err != nil {
		return nil, nil, err
	}
	// get order information
	orderRes := <-srv.orderResp.FindById(ctx, opt.OrderId)
	if orderRes.Error != nil {
		return nil, nil, orderRes.Error
	}
	order = orderRes.Result.(*models.Order)
	if order.OrderNo != opt.OrderNo {
		return nil, nil, errors.New("invalid OrderNo with different")
	}
	if err := order.CanApplyRefund(); err != nil {
		return nil, nil, err
	}
	// valid refund from order
	refund, err := srv.validRefund(order, opt.RefundNo)
	if err != nil {
		return nil, nil, err
	}

	option := srv.buildRefundOption(order, refund)

	response, err := wechat.Pay.Refund(option)
	if err != nil {
		// 参数错误
		return nil, nil, err
	}

	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return nil, nil, err
	}
	if err = session.StartTransaction(); err != nil {
		return nil, nil, err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if response.ReturnCode == "SUCCESS" {
			refund.Status = models.RefundStatusAgreed
			refund.UpdatedAt = time.Now()
			// 接口调用成功
			if response.ResultCode == "SUCCESS" {
				// 退款申请接收成功
				refund.Status = models.RefundStatusRefunding
				// 查询失败记录全部删除
				if err := srv.failedRefundRep.ClearFailedByRefundOn(sessionContext, refund.Id); err != nil {
					session.AbortTransaction(sessionContext)
					log.Printf("删除失败退款记录失败,err:%s", err)
					return err
				}
			} else {
				// 提交业务失败
				log.Printf("refund_no[%s]，提交退款失败！，错误码[%s]，失败原因[%s]", refund.Id, response.ErrCode, response.ErrCodeDes)
				// 记录失败原因
				if err := srv.failedRefundRep.Write(sessionContext, &models.FailedRefund{
					RefundOn:   refund.Id,
					OrderNo:    order.OrderNo,
					ErrCode:    response.ErrCode,
					ErrCodeDes: response.ErrCodeDes,
				}); err != nil {
					session.AbortTransaction(sessionContext)
					log.Printf("保存退款异常记录失败,err:%s", err)
					return err
				}
				err = err2.Err422.F("提交退款失败！，错误码[%s]，失败原因[%s]", response.ErrCode, response.ErrCodeDes)
			}
			saved := <-srv.orderResp.Save(sessionContext, order)
			if saved.Error != nil {
				session.AbortTransaction(sessionContext)
				log.Printf("agree refund,save order error:%s", saved.Error)
				return saved.Error
			}
		}
		session.CommitTransaction(sessionContext)
		return err
	})
	session.EndSession(ctx)
	return refund, order, err
}

// 拒绝退款
func (srv *RefundService) RejectRefund(ctx context.Context, opt *RefundOption) (refund *models.Refund, refundOrder *models.Order, err error) {
	if err := opt.IsValid(); err != nil {
		return nil, nil, err
	}
	// get order information
	orderRes := <-srv.orderResp.FindById(ctx, opt.OrderId)
	if orderRes.Error != nil {
		return nil, nil, orderRes.Error
	}
	order := orderRes.Result.(*models.Order)
	if order.OrderNo != opt.OrderNo {
		return nil, nil, errors.New("invalid OrderNo with different")
	}
	if order.Status != models.OrderStatusPreSend && !order.RefundChannel {
		return nil, nil, errors.New("invalid order status not refund apply")
	}

	refund, err = srv.validRefund(order, opt.RefundNo)
	if err != nil {
		return nil, nil, err
	}

	refund.Status = models.RefundStatusReject
	refund.RejectDesc = opt.Desc
	refund.UpdatedAt = time.Now()

	saved := <-srv.orderResp.Save(ctx, order)
	if saved.Error != nil {
		return nil, nil, saved.Error
	}
	return refund, saved.Result.(*models.Order), nil
}

// 退款回调
func (srv *RefundService) RefundNotify(ctx context.Context, req *http.Request) (order *models.Order, refundNo string, err error) {
	// parse
	notifyReq, err := wechat.Pay.ParseRefundNotifyResult(req)
	if err != nil {
		return nil, "", err
	}
	spew.Dump(notifyReq)
	// deal with
	orderNo := notifyReq.OutTradeNo
	refundNo = notifyReq.OutRefundNo
	order, err = srv.orderResp.FindByOrderNo(ctx, orderNo)
	if err != nil {
		// 订单不存在
		return nil, refundNo, err
	}
	refund, err := order.FindRefund(refundNo)
	if err != nil {
		// 退款单不存在
		return
	}

	if refund.ReturnCode == "SUCCESS" {
		// 已经标记为退款成功
		// 不做处理
		return nil, "", nil
	}

	if refund.ReturnCode == "REFUNDCLOSE" {
		// 当前退款单已经处于退款关闭状态，不做处理
		return nil, "", nil
	}

	refund.ReturnCode = notifyReq.RefundStatus
	if notifyReq.RefundStatus == "SUCCESS" && refund.Status == models.RefundStatusRefunding {
		refund.Status = models.RefundStatusDone
	}
	if notifyReq.RefundStatus == "REFUNDCLOSE" && refund.Status == models.RefundStatusRefunding {
		refund.Status = models.RefundStatusClosed
	}
	refund.UpdatedAt = time.Now()

	// 保存
	saved := <-srv.orderResp.Save(ctx, order)
	if saved.Error != nil {
		// 保存失败
		return nil, "", saved.Error
	}
	savedOrder := saved.Result.(*models.Order)
	return savedOrder, refundNo, nil
}

// 取消退款
func (srv *RefundService) CancelRefund(ctx context.Context, opt *RefundOption, authenticatable auth.Authenticatable, closeChannel bool) (refund *models.Refund, order *models.Order, err error) {
	// 用户取消退款，取消后关闭退款通道
	if err := opt.IsValid(); err != nil {
		return nil, nil, err
	}
	// get order information
	orderRes := <-srv.orderResp.FindById(ctx, opt.OrderId)
	if orderRes.Error != nil {
		return nil, nil, orderRes.Error
	}
	order = orderRes.Result.(*models.Order)
	if order.OrderNo != opt.OrderNo {
		return nil, nil, errors.New("invalid OrderNo with different")
	}
	refund, err = order.FindRefund(opt.RefundNo)
	if err != nil {
		// 退款不存在
		return nil, nil, err
	}
	if !refund.CanCancel() {
		// 当前状态不允许关闭退款
		return nil, nil, err2.Err422.F("当前状态无法关闭退款单")
	}

	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return nil, nil, err
	}
	if err = session.StartTransaction(); err != nil {
		return nil, nil, err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		refund.Status = models.RefundStatusClosed
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
		refund.Canceler = canceler
		refund.UpdatedAt = time.Now()

		// 删除失败退款记录
		if err := srv.failedRefundRep.ClearFailedByRefundOn(sessionContext, refund.Id); err != nil {
			session.AbortTransaction(sessionContext)
			log.Printf("删除异常退款记录失败,err:%s", err)
			return err
		}
		if closeChannel {
			// 关闭退款通道
			order.RefundChannel = false
		}
		// 刷新订单状态
		// order.RefreshStatus()

		// 保存订单
		saved := <-srv.orderResp.Save(sessionContext, order)
		if saved.Error != nil {
			session.AbortTransaction(sessionContext)
			log.Printf("关闭退款单，保存订单失败,err:%s", err)
			return saved.Error
		}
		order = saved.Result.(*models.Order)
		session.CommitTransaction(sessionContext)
		return nil
	})

	session.EndSession(ctx)
	return refund, order, err
}

func (srv *RefundService) buildRefundOption(order *models.Order, refundOrder *models.Refund) (opt *wechat.RefundOption) {
	opt = &wechat.RefundOption{
		OutTradeNo:  refundOrder.OrderNo,
		OutRefundNo: refundOrder.Id,
		TotalFee:    order.ActualAmount,
		RefundFee:   refundOrder.TotalAmount,
		RefundDesc:  refundOrder.RefundDesc,
	}
	// return
	return opt
}

func (srv *RefundService) validRefund(order *models.Order, refundNo string) (refund *models.Refund, err error) {
	refund, err = order.FindRefund(refundNo)
	if err != nil {
		// 退款单不存在
		return nil, err
	}
	if refund.Status == models.RefundStatusApply {
		return refund, nil
	}
	// 不是合法退款状态
	return nil, err2.Err422.F("refund no[%s]，不合法,当前退款单状态[%d],需要状态为[%d]", refundNo, refund.Status, models.RefundStatusApply)
}
