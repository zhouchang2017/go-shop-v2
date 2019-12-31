package resources

import (
	"go-shop-v2/app/vue/pages"
	"go-shop-v2/pkg/support"
	"go-shop-v2/pkg/vue/core"
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

		fx.Invoke(func(vue *core.Vue, resource *Admin) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *core.Vue, resource *Brand) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *core.Vue, resource *Category) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *core.Vue, resource *Product) {
			vue.RegisterResource(resource)
			vue.RegisterPage(pages.NewProductCreatePage())
			vue.RegisterPage(pages.NewProductUpdatePage())
		}),

		fx.Invoke(func(vue *core.Vue, resource *Item) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *core.Vue, resource *Shop) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *core.Vue, resource *Inventory) {
			vue.RegisterResource(resource)
		}),

		fx.Invoke(func(vue *core.Vue, resource *ManualInventoryAction) {
			vue.RegisterResource(resource)
		}),

	)
}
