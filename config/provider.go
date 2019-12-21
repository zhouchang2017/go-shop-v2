package config

import (
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/support"
	"go-shop-v2/pkg/vue"
	"go.uber.org/fx"
)

type configServiceProvider struct {
}

func NewConfigServiceProvider() support.ServiceProvider {
	return &configServiceProvider{}
}

func (c *configServiceProvider) Register(container support.Container) {
	container.Provide(c.mongodbConfig, c.mysqlConfig)
	container.Provide(c.qiniuConfig)
	container.Provide(c.authGuard)
}

func (c *configServiceProvider) Boot() fx.Option {
	return fx.Options(
		fx.Invoke(func(vue *vue.Vue) {
			vue.SetGuard("admin")
		}),
		// 注册七牛http
		fx.Invoke(func(vue *vue.Vue, qiniu *qiniu.Qiniu) {
			vue.RegisterCustomHttpHandler(qiniu.HttpHandle)
		}),
		// 注册事件监听
		fx.Invoke(c.eventRegister()...),
	)
}
