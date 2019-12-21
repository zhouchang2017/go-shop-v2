package models

import (
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
)

var models []interface{}

func register(factory interface{}) {
	models = append(models, factory)
}

type modelServiceProvider struct {
}

func NewModelServiceProvider() support.ServiceProvider {
	return &modelServiceProvider{}
}

func (m *modelServiceProvider) Register(container support.Container) {

	for _,factory:=range models {
		container.Provide(factory)
	}

}

func (m *modelServiceProvider) Boot() fx.Option {
	return nil
}
