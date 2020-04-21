package listeners

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/email"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	mp_subscribe "go-shop-v2/app/mp-subscribe"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/wechat"
)

// 订单支付，通知后台
type OrderPaidNotifyToAdmin struct {
	orderSrv *services.OrderService
}

func (o OrderPaidNotifyToAdmin) Name() string {
	return "订单付款通知"
}

func (o OrderPaidNotifyToAdmin) Make() rabbitmq.Listener {
	return &OrderPaidNotifyToAdmin{orderSrv: services.MakeOrderService()}
}

func (o OrderPaidNotifyToAdmin) Event() rabbitmq.Event {
	return events.OrderPaid{}
}

func (o OrderPaidNotifyToAdmin) OnError(payload []byte, err error) {
	log.Errorf("OrderPaidNotifyToAdmin Error:%s\n", err)
}

// data为订单号
func (o OrderPaidNotifyToAdmin) Handler(data []byte) error {
	orderNo := string(data)
	if orderNo != "" {
		order, err := o.orderSrv.FindByNo(context.Background(), orderNo)
		if err != nil {
			return err
		}
		return o.sendEmailNotifyAdmin(order)
	}
	return nil
}

func (o OrderPaidNotifyToAdmin) sendEmailNotifyAdmin(order *models.Order) error {
	return sendEmail(o.Event(),email.OrderPaidNotify(order))
}

// 订单付款，通知用户
type OrderPaidNotifyToUser struct {
	orderSrv *services.OrderService
}

func (o OrderPaidNotifyToUser) Make() rabbitmq.Listener {
	return &OrderPaidNotifyToUser{
		orderSrv: services.MakeOrderService(),
	}
}

func (o OrderPaidNotifyToUser) Event() rabbitmq.Event {
	return events.OrderPaid{}
}

func (o OrderPaidNotifyToUser) OnError(payload []byte, err error) {
	log.Errorf("OrderPaidNotifyToUser Error:%s\n", err)
}

func (o OrderPaidNotifyToUser) Handler(data []byte) error {
	orderNo := string(data)
	if orderNo != "" {
		order, err := o.orderSrv.FindByNo(context.Background(), orderNo)
		if err != nil {
			return err
		}
		return o.sendWxPush(order)
	}
	return nil
}

func (o OrderPaidNotifyToUser) sendWxPush(order *models.Order) error {
	return wechat.SDK.SendSubscribeMessage(mp_subscribe.NewOrderPaidNotify(order))
}
