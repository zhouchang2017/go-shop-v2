package listeners

import (
	"context"
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

// 用户主动关闭订单
type OrderCloseNotifyToAdmin struct {
	orderSrv *services.OrderService
}

func (o OrderCloseNotifyToAdmin) Make() rabbitmq.Listener {
	return &OrderCloseNotifyToAdmin{orderSrv: services.MakeOrderService()}
}

func (o OrderCloseNotifyToAdmin) Event() rabbitmq.Event {
	return events.OrderClosedByUser{}
}

func (o OrderCloseNotifyToAdmin) OnError(payload []byte, err error) {
	log.Error( payload, err)
}

func (o OrderCloseNotifyToAdmin) Handler(data []byte) error {
	log.Infof("订单付款事件，处理")
	orderNo := string(data)
	spew.Dump(orderNo)
	if orderNo != "" {
		order, err := o.orderSrv.FindByNo(context.Background(), orderNo)
		if err != nil {
			return err
		}
		return o.sendEmailNotifyAdmin(order)
	}
	return nil
}

func (o OrderCloseNotifyToAdmin) sendEmailNotifyAdmin(order *models.Order) error {
	log.Printf("订单关闭通知，发送邮件")
	return email.Send(email.OrderClosedNotify(order, "zhouchangqaz@gmail.com"))
}

// 管理员关闭订单，通知用户
type OrderClosedNotifyToUser struct {
}

func (o OrderClosedNotifyToUser) Make() rabbitmq.Listener {
	return &OrderClosedNotifyToUser{}
}

func (o OrderClosedNotifyToUser) Event() rabbitmq.Event {
	return events.OrderClosedByAdmin{}
}

func (o OrderClosedNotifyToUser) OnError(data []byte, err error) {
	log.Error("OrderClosedNotifyToUser Error:%s\n", err)
}

func (o OrderClosedNotifyToUser) Handler(data []byte) error {
	log.Printf("管理员关闭订单事件，处理")
	var order models.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return err
	}
	return wechat.SDK.SendSubscribeMessage(mp_subscribe.NewOrderClosedNotify(&order))
}
