package message

import (
	"github.com/streadway/amqp"
	"sync"
)

// RabbitMQ 用于管理和维护rabbitmq的对象
type mq struct {
	wg sync.WaitGroup

	channel      *amqp.Channel
	exchangeName string // exchange的名称
	exchangeType string // exchange的类型
	receivers    []Receiver
}
