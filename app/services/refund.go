package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/wechat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type RefundService struct {
	refundResp *repositories.RefundRep
	orderResp  *repositories.OrderRep
}

func NewRefundService(refundResp *repositories.RefundRep, orderResp *repositories.OrderRep) *RefundService {
	return &RefundService{refundResp: refundResp, orderResp: orderResp}
}

type RefundOption struct {
	OrderId string `json:"order_id" form:"order_id" binding:"required"`
	OrderNo string `json:"order_no" form:"order_no" binding:"required"`
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

// 同意退款/进行退款(目前做的整单退，如果要做单个的话需要调整RefundOption添加item)
func (srv *RefundService) AgreeRefund(ctx context.Context, opt *RefundOption) (refundOrder *models.Refund, err error) {
	if err := opt.IsValid(); err != nil {
		return nil, err
	}
	// get order information
	orderRes := <-srv.orderResp.FindById(ctx, opt.OrderId)
	if orderRes.Error != nil {
		return nil, orderRes.Error
	}
	order := orderRes.Result.(*models.Order)
	if order.OrderNo != opt.OrderNo {
		return nil, errors.New("invalid OrderNo with different")
	}
	if order.Status != models.OrderStatusRefundApply {
		return nil, errors.New("invalid order status not refund apply")
	}
	// generate refund order information
	refundOrder = srv.buildRefund(order, order.ActualAmount)
	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return nil, err
	}
	if err := session.StartTransaction(); err != nil {
		return nil, err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// add refund order into db
		created := <-srv.refundResp.Create(sessionContext, refundOrder)
		if created.Error != nil {
			session.AbortTransaction(sessionContext)
			return created.Error
		}
		// update order status into db
		updated := <-srv.orderResp.Update(sessionContext, order.GetID(), bson.M{
			"$set": bson.M{
				"status": models.OrderStatusRefunding,
			},
		})
		if updated.Error != nil {
			session.AbortTransaction(sessionContext)
			return updated.Error
		}
		// refund order
		option := srv.buildRefundOption(refundOrder)
		_, wcErr := wechat.Pay.Refund(option)
		if wcErr != nil {
			session.AbortTransaction(sessionContext)
			return wcErr
		}
		// return
		session.CommitTransaction(sessionContext)
		refundOrder = created.Result.(*models.Refund)
		return nil
	})
	// return
	session.EndSession(ctx)
	return refundOrder, err
}

// 拒绝退款
func (srv *RefundService) RejectRefund(ctx context.Context, opt *RefundOption) (refundOrder *models.Order, err error) {
	if err := opt.IsValid(); err != nil {
		return nil, err
	}
	// get order information
	orderRes := <-srv.orderResp.FindById(ctx, opt.OrderId)
	if orderRes.Error != nil {
		return nil, orderRes.Error
	}
	order := orderRes.Result.(*models.Order)
	if order.OrderNo != opt.OrderNo {
		return nil, errors.New("invalid OrderNo with different")
	}
	if order.Status != models.OrderStatusRefundApply {
		return nil, errors.New("invalid order status not refund apply")
	}
	// update order status to reject
	updated := <-srv.orderResp.Update(ctx, order.GetID(), bson.M{
		"$set": bson.M{
			"status": models.OrderStatusRefundReject,
		},
	})
	if updated.Error != nil {
		return nil, updated.Error
	}
	// return
	return updated.Result.(*models.Order), nil
}

// 退款回调
func (srv *RefundService) RefundNotify(ctx context.Context, req *http.Request) (refundOrderNumber string, err error) {
	// parse
	notifyReq, err := wechat.Pay.ParseRefundNotifyResult(req)
	if err != nil {
		return "", err
	}
	spew.Dump(notifyReq)
	// deal with
	orderNo := notifyReq.OutTradeNo
	refundOrderNo := notifyReq.OutRefundNo
	paymentNo := notifyReq.TransactionId

	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return "", err
	}
	if err = session.StartTransaction(); err != nil {
		return "", err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// get order information
		order, err := srv.orderResp.FindByOrderNo(sessionContext, orderNo)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}
		// get refund order information
		refund, err := srv.refundResp.FindByRefundOrderNo(sessionContext, refundOrderNo)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}
		if order.Status != models.OrderStatusRefunding {
			return fmt.Errorf("order %s status %d not refunding", order.OrderNo, order.Status)
		}
		if refund.PaymentStatus == "SUCCESS" {
			// 标记为已成功
			return nil
		}
		// check result
		if notifyReq.RefundStatus != "SUCCESS" {
			// 通知结果支付失败
			updateRes := <-srv.refundResp.Update(sessionContext, refund.GetID(), bson.M{
				"$set": bson.M{
					"payment_status": notifyReq.RefundStatus,
					"payment_no":     paymentNo,
				},
			})
			if updateRes.Error != nil {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("update refund order-%s to failed status failed %s", refundOrderNo, updateRes.Error)
			}
			// return
			session.CommitTransaction(sessionContext)
			return nil
		}
		// deal with if success
		// valid money
		if uint64(notifyReq.TotalFee) != refund.TotalAmount {
			session.AbortTransaction(sessionContext)
			return fmt.Errorf("notify total fee %d is different with refund order total amount %d", notifyReq.TotalFee, refund.TotalAmount)
		}
		// 支付成功
		// update refund order
		updateRefund := <-srv.refundResp.Update(sessionContext, refund.GetID(), bson.M{
			"$set": bson.M{
				"payment_status": notifyReq.RefundStatus,
				"payment_no":     paymentNo,
			},
		})
		if updateRefund.Error != nil {
			session.AbortTransaction(sessionContext)
			return fmt.Errorf("update refund order-%s to success status failed %s", refundOrderNo, updateRefund.Error)
		}
		// update order status
		updateOrder := <-srv.orderResp.Update(sessionContext, order.GetID(), bson.M{
			"$set": bson.M{
				"status": models.OrderStatusRefundDone,
			},
		})
		if updateOrder.Error != nil {
			session.AbortTransaction(sessionContext)
			return fmt.Errorf("update order-%s to refund success status failed %s", orderNo, updateOrder.Error)
		}
		// return
		session.CommitTransaction(sessionContext)
		refundOrderNumber = refundOrderNo // 新标记订单为完成状态，返回退款订单号
		return nil
	})

	session.EndSession(ctx)
	// return
	return refundOrderNumber, err
}

func (srv *RefundService) buildRefundOption(refundOrder *models.Refund) (opt *wechat.RefundOption) {
	opt = &wechat.RefundOption{
		OutTradeNo:  refundOrder.OrderNo,
		OutRefundNo: refundOrder.RefundOrderNo,
		TotalFee:    refundOrder.TotalAmount,
		RefundFee:   refundOrder.TotalAmount,
		RefundDesc:  "整单退款",
	}
	// return
	return opt
}

func (srv *RefundService) buildRefund(order *models.Order, refundFee uint64) (refund *models.Refund) {
	refundItems := make([]*models.RefundItem, 0)
	for _, item := range order.OrderItems {
		refundItem := &models.RefundItem{
			ItemId: item.Item.Id,
			Qty:    item.Count,
			Amount: uint64(item.Price),
		}
		refundItems = append(refundItems, refundItem)
	}
	refundOrder := &models.Refund{
		OrderNo:       order.OrderNo,
		RefundOrderNo: utils.RandomRefundOrderNo(""),
		PaymentNo:     "",
		TotalAmount:   refundFee,
		Items:         refundItems,
		PaymentStatus: "",
	}
	return refundOrder
}
