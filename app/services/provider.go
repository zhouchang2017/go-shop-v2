package services

import (
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
)

var services []interface{}

func register(factory interface{}) {
	services = append(services, factory)
}

type serviceServiceProvider struct {
}

func NewServiceServiceProvider() support.ServiceProvider {
	return &serviceServiceProvider{}
}

func (s *serviceServiceProvider) Register(container support.Container) {
	for _, factory := range services {
		container.Provide(factory)
	}
}

func (s *serviceServiceProvider) Boot() fx.Option {
	return nil
}
