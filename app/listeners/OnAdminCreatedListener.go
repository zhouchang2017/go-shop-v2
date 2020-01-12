package listeners

import (
	"context"
	"encoding/json"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/message"
	"log"
)

// 用户创建 同步门店关联的用户
type AdminCreatedSyncAssociatedShop struct {
	shopSrv *services.ShopService
}

func NewAdminCreatedSyncAssociatedShop() *AdminCreatedSyncAssociatedShop {
	return &AdminCreatedSyncAssociatedShop{shopSrv:services.MakeShopService()}
}

// 对应的触发事件
func (AdminCreatedSyncAssociatedShop) Event() message.Event {
	return events.AdminCreated{}
}

// 队列名称
func (AdminCreatedSyncAssociatedShop) QueueName() string {
	return "AdminCreatedSyncAssociatedShop"
}

// 错误处理
func (AdminCreatedSyncAssociatedShop) OnError(err error) {
	log.Printf("AdminCreatedSyncAssociatedShop Error:%s\n", err)
}

// 处理逻辑
func (this AdminCreatedSyncAssociatedShop) Handler(data []byte) error {
	var admin models.Admin
	if err := json.Unmarshal(data, &admin); err != nil {
		return err
	}
	return this.shopSrv.SyncAssociatedMembers(context.Background(), &admin)
}
