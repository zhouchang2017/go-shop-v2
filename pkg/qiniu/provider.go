package qiniu

import (
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
)

type qiniuServiceProvider struct {
}

func (q *qiniuServiceProvider) Register(container support.Container) {
	container.Provide(NewQiniu)
}

func (q *qiniuServiceProvider) Boot() fx.Option {
	return nil
}

func NewQiniuServiceProvider() support.ServiceProvider {
	return &qiniuServiceProvider{}
}
