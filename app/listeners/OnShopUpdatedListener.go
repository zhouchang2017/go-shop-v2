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

// 门店更新 同步用户关联的门店
type ShopUpdatedSyncAssociatedAdmin struct {
	adminSrv *services.AdminService
}

func NewShopUpdatedSyncAssociatedAdmin() *ShopUpdatedSyncAssociatedAdmin {
	return &ShopUpdatedSyncAssociatedAdmin{adminSrv:services.MakeAdminService()}
}

// 对应的触发事件
func (ShopUpdatedSyncAssociatedAdmin) Event() message.Event {
	return events.ShopUpdated{}
}

// 队列名称
func (ShopUpdatedSyncAssociatedAdmin) QueueName() string {
	return "ShopUpdatedSyncAssociatedAdmin"
}

// 错误处理
func (ShopUpdatedSyncAssociatedAdmin) OnError(err error) {
	log.Printf("ShopUpdatedSyncAssociatedAdmin Error:%s\n", err)
}

// 处理逻辑
func (this ShopUpdatedSyncAssociatedAdmin) Handler(data []byte) error {
	var shop models.Shop
	err := json.Unmarshal(data, &shop)
	if err != nil {
		return err
	}
	return this.adminSrv.SyncAssociatedShop(context.Background(), &shop)
}

