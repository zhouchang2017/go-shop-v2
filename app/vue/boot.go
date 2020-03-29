package vue

import (
	"go-shop-v2/app/vue/charts"
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
		pages2.NewProductCreatePage(),   // 产品创建页
		pages2.NewProductUpdatePage(),   // 产品更新页
		pages2.NewPromotionCreatePage(), // 促销计划创建页
		pages2.NewPromotionUpdatePage(), // 促销计划更新页

		pages2.NewOrderShipmentPage(),  // 订单发货页面
		pages2.NewOrderLogisticsPage(), // 物流详情页面

		pages2.NewOrderItemAggregatePage(), // 订单明细聚合
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

	// dashboard页面cards注册
	vue.RegisterCard(charts.NewNewUserValue())
	vue.RegisterCard(charts.NewNewOrderValue())
	vue.RegisterCard(charts.NewNewPaymentValue())
	vue.RegisterCard(charts.NewCountOrderPrePayValue())
	vue.RegisterCard(charts.NewCountOrderPreSendValue())
	vue.RegisterCard(charts.NewPaymentLine())
}
