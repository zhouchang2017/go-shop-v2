package listeners

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/email"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/rabbitmq"
)

type OnOrderCreatedListener struct {
}

func (o OnOrderCreatedListener) Name() string {
	return "新订单通知"
}

func (o OnOrderCreatedListener) Make() rabbitmq.Listener {
	return &OnOrderCreatedListener{}
}

func (o OnOrderCreatedListener) OnError(payload []byte, err error) {
	log.Errorf("OnOrderCreatedListener Error:%s\n", err)
}

func NewOnOrderCreatedListener() *OnOrderCreatedListener {
	return &OnOrderCreatedListener{}
}

// 对应的触发事件
func (o OnOrderCreatedListener) Event() rabbitmq.Event {
	return events.OrderCreated{}
}

// 处理逻辑
func (o OnOrderCreatedListener) Handler(data []byte) error {
	var order models.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return err
	}
	return o.sendEmailNotifyAdmin(&order)
}

// 发送邮件给管理员
func (o OnOrderCreatedListener) sendEmailNotifyAdmin(order *models.Order) error {
	return sendEmail(o.Event(),email.OrderCreatedNotify(order))
}
