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
		//resources2.NewInventoryResource(), // 库存系统
		resources2.NewItemResource(),
		//resources2.NewInventoryActionResource(), // 库存系统操作
		resources2.NewProductResource(),
		resources2.NewShopResource(),
		resources2.NewOrderResource(),
		//resources2.NewInventoryLogResource(), // 库存系统日志
		resources2.NewArticleResource(),
		resources2.NewTopicResource(),
		resources2.NewPromotionResource(),
		resources2.NewPromotionItemResource(),
	}
}

func registerPages() {
	pages = []contracts.Page{
		//pages2.NewInventoryAggregatePage(), // 聚合库存自定义页面
		//pages2.NewManualInventoryCreatePage(), // 入库\出库添加页面
		//pages2.NewManualInventoryUpdatePage(), // 入库\出库更新页面
		pages2.NewProductCreatePage(),
		pages2.NewProductUpdatePage(),
		pages2.NewPromotionCreatePage(),
		pages2.NewPromotionUpdatePage(),
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
