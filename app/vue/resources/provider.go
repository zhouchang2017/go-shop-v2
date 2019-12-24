package resources

import (
	"go-shop-v2/pkg/support"
	"go-shop-v2/pkg/vue"
	"go.uber.org/fx"
)

var resources []interface{}

func register(factory interface{}) {
	resources = append(resources, factory)
}

var resourceList list

type list struct {
	fx.In
	Resources []interface{} `group:"resource"`
}

type vueResourceServiceProvider struct {
}

func NewVueResourceServiceProvider() support.ServiceProvider {
	return &vueResourceServiceProvider{}
}

func (v *vueResourceServiceProvider) Register(container support.Container) {
	for _, factory := range resources {
		container.Provide(factory)
	}
}

func (v *vueResourceServiceProvider) Boot() fx.Option {

	return fx.Options(

		fx.Invoke(func(vue *vue.Vue, resource *Admin) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *vue.Vue, resource *Brand) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *vue.Vue, resource *Category) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *vue.Vue, resource *Product) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *vue.Vue, resource *Item) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *vue.Vue, resource *Shop) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *vue.Vue, resource *Inventory) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *vue.Vue, resource *ManualInventoryAction) {
			vue.RegisterResource(resource)
		}),

	)
}
