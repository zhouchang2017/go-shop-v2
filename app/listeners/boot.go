package listeners

import (
	"go-shop-v2/pkg/message"
)

var listeners []func() message.Listener

func Boot(mq *message.RabbitMQ) {
	for _, factory := range listeners {
		mq.Register(factory())
	}
}

func register(factory func() message.Listener) {
	listeners = append(listeners, factory)
}
