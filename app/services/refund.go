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
	"time"
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

// refund todo: 目前做的整单退，如果要做单个的话需要调整RefundOption添加item
func (srv *RefundService) Refund(ctx context.Context, opt *RefundOption) (*models.Refund, error) {
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
	if order.Status != models.OrderStatusPreSend {
		return nil, errors.New("invalid order status not presend")
	}
	// generate refund order information
	refundOrder := srv.buildRefund(order, order.ActualAmount)
	// add refund order into db
	created := <-srv.refundResp.Create(ctx, refundOrder)
	if created.Error != nil {
		return nil, created.Error
	}
	// refund order
	option := srv.buildRefundOption(order, refundOrder.RefundOrderNo)
	_, wcErr := wechat.Pay.Refund(option)
	if wcErr != nil {
		return nil, wcErr
	}
	return created.Result.(*models.Refund), nil
}

// 退款回调
func (srv *RefundService) RefundNotify(ctx context.Context, req *http.Request) (orderNumber string, err error) {

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
		// 从payments中查询order对应的payment
		payment, err := srv.paymentRep.FindByOrderId(sessionContext, orderNo)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}
		if payment.PaymentAt != nil && order.Status == models.OrderStatusPaid {
			// 已经标记为支付！！
			return nil
		}
		payment.PaymentNo = paymentNo

		// check result
		if notifyReq.ReturnCode == "SUCCESS" {
			// if failed
			if notifyReq.ResultCode != "SUCCESS" {
				// 通知结果支付失败
				// update payment
				order.Payment = payment.ToAssociated()
				order.StatusToFailed()

				updateRes := <-srv.orderResp.Update(sessionContext, order.GetID(), bson.M{
					"$set": bson.M{
						"status":  models.OrderStatusFailed,
						"payment": payment.ToAssociated(),
					},
				})
				if updateRes.Error != nil {
					session.AbortTransaction(sessionContext)
					return fmt.Errorf("update order-%s to failed status failed %s", orderNo, updateRes.Error)
				}
				// return
				//log.Println("update order xxx to failed status success")
				return nil
			}
			// deal with if success
			// valid money
			if uint64(notifyReq.TotalFee) != order.ActualAmount {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("notify total fee %d is different with order actucal amount %d", notifyReq.TotalFee, order.ActualAmount)
			}
			// 支付成功
			// update payment
			payment.SetPaymentAt(time.Now())
			saved := <-srv.paymentRep.Save(sessionContext, payment)
			if saved.Error != nil {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("update order[%s] payment to paided failed %s", orderNo, saved.Error)
			}
			updatePayment := saved.Result.(*models.Payment)
			// update order
			updateRes := <-srv.orderResp.Update(sessionContext, order.GetID(), bson.M{
				"$set": bson.M{
					"status":  models.OrderStatusPreSend,
					"payment": updatePayment.ToAssociated(),
				},
			})
			if updateRes.Error != nil {
				session.AbortTransaction(sessionContext)
				return fmt.Errorf("update order-%s to success status failed %s", orderNo, updateRes.Error)
			}
			// todo: add other case to do and remember use transaction if involving other table in db

			// return
			//log.Println("update order xxx to success status success")
			session.CommitTransaction(sessionContext)

			// 新标记订单为支付状态，返回订单号
			orderNumber = orderNo
			return nil
		} else {
			session.AbortTransaction(sessionContext)
			return errors.New("return code not SUCCESS")
		}
	})

	session.EndSession(ctx)
	// return
	return orderNumber, err
}

func (srv *RefundService) buildRefundOption(order *models.Order, refundOrderNo string) (opt *wechat.RefundOption) {
	opt = &wechat.RefundOption{
		OutTradeNo:  order.OrderNo,
		OutRefundNo: refundOrderNo,
		TotalFee:    order.ActualAmount,
		RefundFee:   order.ActualAmount,
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
		Status:        models.RefundStatusPending,
	}
	return refundOrder
}
