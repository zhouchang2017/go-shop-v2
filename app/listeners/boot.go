package listeners

import (
	"go-shop-v2/pkg/message"
)

func Boot(mq *message.RabbitMQ) {
	mq.Register(NewAdminCreatedSyncAssociatedShop())
	mq.Register(NewAdminUpdatedSyncAssociatedShop())
	mq.Register(NewShopCreatedSyncAssociatedAdmin())
	mq.Register(NewShopUpdatedSyncAssociatedAdmin())
	//mq.Register(NewTimeOutCloseInventoryAction())
}

// 前端事件注册
func FrontEndBoot(mq *message.RabbitMQ) {
	mq.Register(NewOnOrderCreatedListener()) // 新订单
	mq.Register(NewOnOrderPaidListener())    // 订单已付款
}
