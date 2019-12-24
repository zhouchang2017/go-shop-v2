package policies

import (
	"go-shop-v2/pkg/support"
	"go-shop-v2/pkg/vue"
	"go.uber.org/fx"
)

var policies []interface{}

func register(factory interface{}) {
	policies = append(policies, factory)
}

type PolicyServiceProvider struct {
}

func NewPolicyServiceProvider() support.ServiceProvider {
	return &PolicyServiceProvider{}
}

func (*PolicyServiceProvider) Register(container support.Container) {
	for _, factory := range policies {
		container.Provide(factory)
	}
}

func (*PolicyServiceProvider) Boot() fx.Option {
	return fx.Options(
		fx.Invoke(func(vue *vue.Vue, policy *InventoryPolicy) {
			vue.RegisterPolice(policy)
		}),
		fx.Invoke(func(vue *vue.Vue, policy *ManualInventoryActionPolicy) {
			vue.RegisterPolice(policy)
		}),
	)
}
