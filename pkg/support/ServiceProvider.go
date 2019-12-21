package support

import (
	"go.uber.org/fx"
)

type ServiceProvider interface {
	// 注册
	Register(container Container)
	// 启动
	Boot() fx.Option
}

type Container interface {
	Provide(constructors ...interface{})
}

type fxOptions struct {
	opts []fx.Option
}

func (o *fxOptions) Provide(constructors ...interface{}) {
	o.opts = append(o.opts, fx.Provide(constructors...))
}

func (o *fxOptions) merge(opt fx.Option) {
	o.opts = append(o.opts, opt)
}

func Run(factories ...func() ServiceProvider) {
	options := &fxOptions{}
	var serviceProviders []ServiceProvider
	for _, factory := range factories {
		provider := factory()
		serviceProviders = append(serviceProviders, provider)
		// register
		provider.Register(options)
	}

	for _, provider := range serviceProviders {
		// boot
		if opt := provider.Boot(); opt != nil {
			options.merge(opt)
		}
	}

	app := fx.New(fx.Options(options.opts...))

	app.Run()
}
