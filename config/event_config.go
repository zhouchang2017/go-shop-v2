package config

import (
	"go-shop-v2/app/events"
	"go-shop-v2/app/listeners"
	"go-shop-v2/pkg/event"
)

var eventRegisters []interface{} = []interface{}{
	func(syncAdminAssociatedShop *listeners.SyncAdminAssociatedShop) {
		// 创建门店事件
		event.AddListen(events.ShopCreated{}, syncAdminAssociatedShop)
		// 更新门店事件
		event.AddListen(events.ShopUpdated{}, syncAdminAssociatedShop)
	},
	func(syncShopAssociatedAdmin *listeners.SyncShopAssociatedAdmin) {
		// 创建后台用户事件
		event.AddListen(events.AdminCreated{}, syncShopAssociatedAdmin)
		// 更新后台用户事件
		event.AddListen(events.AdminUpdated{}, syncShopAssociatedAdmin)
	},
}

func (c *configServiceProvider) eventRegister() []interface{} {
	return eventRegisters
}
