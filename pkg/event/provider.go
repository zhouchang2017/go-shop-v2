package event

import (
	"context"
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
)

type eventServiceProvider struct {
	cancel context.CancelFunc
}

func (e *eventServiceProvider) Register(container support.Container) {
	container.Provide(NewBus)
}

func (e *eventServiceProvider) Boot() fx.Option {

	fx.Invoke()

	return fx.Invoke(e.start)
}

func NewEventServiceProvider() support.ServiceProvider {
	return &eventServiceProvider{}
}

func (e *eventServiceProvider) start(lifecycle fx.Lifecycle, bus *Bus) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ctx2, cancel := context.WithCancel(context.Background())
			e.cancel = cancel
			bus.Run(ctx2)
			return nil
		},
		OnStop: func(i context.Context) error {
			e.cancel()
			return nil
		},
	})
}
