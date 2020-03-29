package listeners

import (
	"go-shop-v2/pkg/rabbitmq"
)

func Boot(mq *rabbitmq.RabbitMQ) {
	mq.Register(OnOrderCreatedListener{}.Make())  // 订单创建事件处理
	mq.Register(OrderPaidNotifyToAdmin{}.Make())  // 订单付款通知管理员
	mq.Register(OrderPaidNotifyToUser{}.Make())   // 订单付款通知客户
	mq.Register(OrderCloseNotifyToAdmin{}.Make()) // 买家关闭订单通知管理员
	mq.Register(OrderClosedNotifyToUser{}.Make()) // 管理员关闭订单通知用户
	mq.Register(OnOrderTimeOutListener{}.Make())  // 订单超时自动关闭
}
