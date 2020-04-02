package listeners

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/events"
	mp_subscribe "go-shop-v2/app/mp-subscribe"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/wechat"
)

type OnOrderTimeOutListener struct {
	orderSrv *services.OrderService
}

func (o OnOrderTimeOutListener) Make() rabbitmq.Listener {
	return &OnOrderTimeOutListener{orderSrv: services.MakeOrderService()}
}

func (o OnOrderTimeOutListener) Event() rabbitmq.Event {
	return events.OrderTimeOut{}
}

func (o OnOrderTimeOutListener) OnError(payload []byte, err error) {
	log.Errorf("OnOrderTimeOutListener Error:%s\n", err)
}

func (o OnOrderTimeOutListener) Handler(data []byte) error {
	log.Printf("订单超时自动关闭，处理")
	orderNo := string(data)
	order, err := o.orderSrv.FindByNo(context.Background(), orderNo)
	if err != nil {
		return err
	}
	if order.StatusIsFailed() {
		// 订单已经被关闭,不做处理
		return nil
	}
	if order.StatusIsFailed() {
		// 已关闭订单，不做处理
		return nil
	}
	// 关闭订单
	updatedOrder, err := o.orderSrv.Cancel(context.Background(), order, "")
	if err != nil {
		return err
	}
	// 小程序推送
	return wechat.SDK.SendSubscribeMessage(mp_subscribe.NewOrderClosedNotify(updatedOrder))
}
