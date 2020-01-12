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

// 用户更新 同步门店关联的用户
type AdminUpdatedSyncAssociatedShop struct {
	shopSrv *services.ShopService
}

func NewAdminUpdatedSyncAssociatedShop() *AdminUpdatedSyncAssociatedShop {
	return &AdminUpdatedSyncAssociatedShop{shopSrv:services.MakeShopService()}
}

// 对应的触发事件
func (AdminUpdatedSyncAssociatedShop) Event() message.Event {
	return events.AdminUpdated{}
}

// 队列名称
func (AdminUpdatedSyncAssociatedShop) QueueName() string {
	return "AdminUpdatedSyncAssociatedShop"
}

// 错误处理
func (AdminUpdatedSyncAssociatedShop) OnError(err error) {
	log.Printf("AdminUpdatedSyncAssociatedShop Error:%s\n", err)
}

// 处理逻辑
func (this AdminUpdatedSyncAssociatedShop) Handler(data []byte) error {
	var admin models.Admin
	if err := json.Unmarshal(data, &admin); err != nil {
		return err
	}
	return this.shopSrv.SyncAssociatedMembers(context.Background(), &admin)
}

