package listeners

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/email"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	mp_subscribe "go-shop-v2/app/mp-subscribe"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/wechat"
)

type OnOrderApplyRefundListener struct {
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
	var event events.OrderApplyRefund
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}
	return o.sendEmailNotifyAdmin(event.Order, event.RefundId)
}

func (o OnOrderApplyRefundListener) sendEmailNotifyAdmin(order *models.Order, refundId string) error {
	log.Printf("订单申请退款通知，发送邮件")
	return email.Send(email.OrderApplyRefundNotify(order, refundId, "zhouchangqaz@gmail.com"))
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
	var event events.OrderRefundChange
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}
	return o.sendWechatNotifyUser(event.Order, event.RefundId)
}

func (o OnOrderRefundChangeListener) sendWechatNotifyUser(order *models.Order, refundId string) error {
	log.Printf("退款订单状态变动，微信推送给用户")
	return wechat.SDK.SendSubscribeMessage(mp_subscribe.NewOrderRefundChangeNotify(order, refundId))
}

type OnOrderRefundCancelByUserListener struct {
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
	var event events.OrderRefundCancelByUser
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}
	return o.sendEmailNotifyAdmin(event.Order, event.RefundId)
}

func (o OnOrderRefundCancelByUserListener) sendEmailNotifyAdmin(order *models.Order, refundId string) error {
	log.Printf("订单申请退款通知，发送邮件")
	return email.Send(email.OrderApplyRefundNotify(order, refundId, "zhouchangqaz@gmail.com"))
}

type OnOrderRefundDoneListener struct {
	orderSrv *services.OrderService
}

func (o OnOrderRefundDoneListener) Make() rabbitmq.Listener {
	return &OnOrderRefundDoneListener{orderSrv: services.MakeOrderService()}
}

func (o OnOrderRefundDoneListener) Event() rabbitmq.Event {
	return events.OrderRefundDone{}
}

func (o OnOrderRefundDoneListener) OnError(payload []byte, err error) {
	log.Errorf("OnOrderRefundDoneListener Error:%s\n", err)
}

func (o OnOrderRefundDoneListener) Handler(data []byte) error {
	orderId := string(data)
	order, err := o.orderSrv.FindById(context.Background(), orderId)
	if err != nil {
		return err
	}
	order.RefreshStatus()
	_, err = o.orderSrv.Save(context.Background(), order)
	return err
}
