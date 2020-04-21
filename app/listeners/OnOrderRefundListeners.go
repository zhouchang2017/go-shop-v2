package listeners

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/email"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	mp_subscribe "go-shop-v2/app/mp-subscribe"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/wechat"
)

// 退款状态变化监听
type OrRefundChangeListener struct {
	orderSrv  *services.OrderService
	refundSrv *services.RefundService
}

func (o OrRefundChangeListener) Make() rabbitmq.Listener {
	return &OrRefundChangeListener{orderSrv: services.MakeOrderService(), refundSrv: services.MakeRefundService()}
}

func (o OrRefundChangeListener) Event() rabbitmq.Event {
	return events.OrderRefundChange{}
}

func (o OrRefundChangeListener) OnError(payload []byte, err error) {
	log.Errorf("OrRefundChangeListener Error:%s\n", err)
}

func (o OrRefundChangeListener) Handler(data []byte) error {
	var refund models.Refund
	if err := json.Unmarshal(data, &refund); err != nil {
		return err
	}
	spew.Dump(refund)
	//o.refundSrv.FindRefundByOrderId(context.Background(),refund.OrderId)
	//panic("implement me")
	return nil
}

type OnOrderApplyRefundListener struct {
}

func (o OnOrderApplyRefundListener) Name() string {
	return "订单申请退款通知"
}

func (o OnOrderApplyRefundListener) Make() rabbitmq.Listener {
	return &OnOrderApplyRefundListener{}
}

func (o OnOrderApplyRefundListener) Event() rabbitmq.Event {
	return events.OrderApplyRefund{}
}

func (o OnOrderApplyRefundListener) OnError(payload []byte, err error) {
	log.Errorf("OnOrderApplyRefundListener Error:%s\n", err)
}

func (o OnOrderApplyRefundListener) Handler(data []byte) error {
	var refund models.Refund
	if err := json.Unmarshal(data, &refund); err != nil {
		return err
	}
	return o.sendEmailNotifyAdmin(&refund)
}

func (o OnOrderApplyRefundListener) sendEmailNotifyAdmin(refund *models.Refund) error {
	log.Printf("订单申请退款通知，发送邮件")
	return sendEmail(o.Event(), email.OrderApplyRefundNotify(refund))
}

type OnOrderRefundChangeListener struct {
}

func (o OnOrderRefundChangeListener) Make() rabbitmq.Listener {
	return &OnOrderRefundChangeListener{}
}

func (o OnOrderRefundChangeListener) Event() rabbitmq.Event {
	return events.OrderRefundChange{}
}

func (o OnOrderRefundChangeListener) OnError(payload []byte, err error) {
	log.Errorf("OnOrderRefundChangeListener Error:%s\n", err)
}

func (o OnOrderRefundChangeListener) Handler(data []byte) error {
	var refund models.Refund
	if err := json.Unmarshal(data, &refund); err != nil {
		return err
	}
	return o.sendWechatNotifyUser(&refund)
}

func (o OnOrderRefundChangeListener) sendWechatNotifyUser(refund *models.Refund) error {
	log.Printf("退款订单状态变动，微信推送给用户")
	return wechat.SDK.SendSubscribeMessage(mp_subscribe.NewOrderRefundChangeNotify(refund))
}

type OnOrderRefundCancelByUserListener struct {
}

func (o OnOrderRefundCancelByUserListener) Name() string {
	return "用户取消退款通知"
}

func (o OnOrderRefundCancelByUserListener) Make() rabbitmq.Listener {
	return &OnOrderRefundCancelByUserListener{}
}

func (o OnOrderRefundCancelByUserListener) Event() rabbitmq.Event {
	return events.OrderRefundCancelByUser{}
}

func (o OnOrderRefundCancelByUserListener) OnError(payload []byte, err error) {
	log.Errorf("OnOrderRefundCancelByUserListener Error:%s\n", err)
}

func (o OnOrderRefundCancelByUserListener) Handler(data []byte) error {
	var refund models.Refund
	if err := json.Unmarshal(data, &refund); err != nil {
		return err
	}
	return o.sendEmailNotifyAdmin(&refund)
}

func (o OnOrderRefundCancelByUserListener) sendEmailNotifyAdmin(refund *models.Refund) error {
	log.Printf("订单申请退款通知，发送邮件")
	return sendEmail(o.Event(), email.OrderApplyRefundNotify(refund))
}
