package listeners

import (
	"encoding/json"
	"go-shop-v2/app/email"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/message"
	"log"
)

type OnOrderCreatedListener struct {
}

func NewOnOrderCreatedListener() *OnOrderCreatedListener {
	return &OnOrderCreatedListener{}
}

// 对应的触发事件
func (o OnOrderCreatedListener) Event() message.Event {
	return events.OrderCreated{}
}

// 队列名称
func (o OnOrderCreatedListener) QueueName() string {
	return "OnOrderCreatedListener"
}

// 错误处理
func (o OnOrderCreatedListener) OnError(err error) {
	log.Printf("OnOrderCreatedListener Error:%s\n", err)
}

// 处理逻辑
func (o OnOrderCreatedListener) Handler(data []byte) error {
	log.Printf("新订单事件，处理")
	var order models.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return err
	}
	return o.sendEmailNotifyAdmin(&order)
}

// 发送邮件给管理员
func (o OnOrderCreatedListener) sendEmailNotifyAdmin(order *models.Order) error {
	log.Printf("新订单通知，发送邮件")
	return email.Send(email.OrderCreatedNotify(order, "zhouchangqaz@gmail.com"))
}
