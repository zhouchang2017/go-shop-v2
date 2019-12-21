package repositories

import (
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
)

var repositories []interface{}

func register(factory interface{}) {
	repositories = append(repositories, factory)
}


type repositoryServiceProvider struct {
}

func NewRepositoryServiceProvider() support.ServiceProvider {
	return &repositoryServiceProvider{}
}

func (r *repositoryServiceProvider) Register(container support.Container) {
	for _, factory := range repositories {
		container.Provide(factory)
	}
}

func (r *repositoryServiceProvider) Boot() fx.Option {
	return nil
}
