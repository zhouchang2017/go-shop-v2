package services

import "go-shop-v2/pkg/rabbitmq"

// 消息队列处理
type EventService struct {
	listeners []rabbitmq.Listener
}
