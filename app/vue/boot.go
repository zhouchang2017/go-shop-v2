package vue

import (
	pages2 "go-shop-v2/app/vue/pages"
	resources2 "go-shop-v2/app/vue/resources"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
)

var resources []contracts.Resource

var pages []contracts.Page

func registerResources() {
	resources = []contracts.Resource{
		resources2.NewAdminResource(),
		resources2.NewBrandResource(),
		resources2.NewCategoryResource(),
		resources2.NewInventoryResource(),
		resources2.NewItemResource(),
		resources2.NewInventoryActionResource(),
		resources2.NewProductResource(),
		resources2.NewShopResource(),
		resources2.NewOrderResource(),
	}
}

func registerPages() {
	pages = []contracts.Page{
		pages2.NewInventoryAggregatePage(),
		pages2.NewManualInventoryCreatePage(),
		pages2.NewManualInventoryUpdatePage(),
		pages2.NewProductCreatePage(),
		pages2.NewProductUpdatePage(),
	}
}

func Boot(vue *core.Vue) {
	registerResources()
	// 注册资源
	for _, resource := range resources {
		vue.RegisterResource(resource)
	}

	registerPages()
	// 注册页面
	for _, page := range pages {
		vue.RegisterPage(page)
	}
}
