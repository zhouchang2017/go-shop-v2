package listeners

import (
	"go-shop-v2/pkg/message"
)

func Boot(mq *message.RabbitMQ) {
	mq.Register(NewAdminCreatedSyncAssociatedShop())
	mq.Register(NewAdminUpdatedSyncAssociatedShop())
	mq.Register(NewShopCreatedSyncAssociatedAdmin())
	mq.Register(NewShopUpdatedSyncAssociatedAdmin())
	mq.Register(NewTimeOutCloseInventoryAction())
}
