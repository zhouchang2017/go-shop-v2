package listeners

import (
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
)

var listeners []interface{}

func register(factory interface{}) {
	listeners = append(listeners, factory)
}

type ListenerServiceProvider struct {
}

func NewListenerServiceProvider() support.ServiceProvider {
	return &ListenerServiceProvider{}
}

func (l *ListenerServiceProvider) Register(container support.Container) {
	for _, factory := range listeners {
		container.Provide(factory)
	}
}

func (l *ListenerServiceProvider) Boot() fx.Option {
	return nil
}
