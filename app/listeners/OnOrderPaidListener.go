package listeners

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/email"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/message"
	"log"
)

type OnOrderPaidListener struct {
	orderSrv *services.OrderService
}

func NewOnOrderPaidListener() *OnOrderPaidListener {
	return &OnOrderPaidListener{
		orderSrv: services.MakeOrderService(),
	}
}

func (o OnOrderPaidListener) Event() message.Event {
	return events.OrderPaid{}
}

func (o OnOrderPaidListener) QueueName() string {
	return "OnOrderPaidListener"
}

func (o OnOrderPaidListener) OnError(err error) {
	log.Printf("OnOrderPaidListener Error:%s\n", err)
}

// data为订单号
func (o OnOrderPaidListener) Handler(data []byte) error {
	log.Printf("订单付款事件，处理")
	orderNo := string(data)
	spew.Dump(orderNo)
	if orderNo != "" {
		order, err := o.orderSrv.FindByNo(context.Background(), orderNo)
		spew.Dump(order)
		if err != nil {
			return err
		}
		return o.sendEmailNotifyAdmin(order)
	}
	return nil
}

func (o OnOrderPaidListener) sendEmailNotifyAdmin(order *models.Order) error {
	log.Printf("订单付款通知，发送邮件")
	return email.Send(email.OrderPaidNotify(order, "zhouchangqaz@gmail.com"))
}
