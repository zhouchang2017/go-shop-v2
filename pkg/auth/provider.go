package auth

import (
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
)

func NewAuthServiceProvider() support.ServiceProvider {
	return &authServiceProvider{}
}

type authServiceProvider struct {
	manager *AuthManager
}

func (a *authServiceProvider) Register(container support.Container) {
	a.manager = NewAuth()
}

func (a *authServiceProvider) Boot() fx.Option {
	return fx.Options(fx.Invoke(a.manager.Register))
}
