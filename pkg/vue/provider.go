package vue

import (
	"context"
	"go-shop-v2/pkg/support"
	"go-shop-v2/pkg/vue/core"
	"go.uber.org/fx"
)

const port = 8083

type vueServiceProvider struct {
}

func NewVueServiceProvider() support.ServiceProvider {
	return &vueServiceProvider{}
}

func (v *vueServiceProvider) Register(container support.Container) {
	container.Provide(func() *core.Vue {
		return core.New(port)
	})
}

func (v *vueServiceProvider) Boot() fx.Option {
	return fx.Invoke(v.start)
}

func (v *vueServiceProvider) start(lifecycle fx.Lifecycle, vue *core.Vue) {

	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			return vue.Run()
		},
		OnStop: func(i context.Context) error {
			return vue.Shutdown(i)
		},
	})
}
