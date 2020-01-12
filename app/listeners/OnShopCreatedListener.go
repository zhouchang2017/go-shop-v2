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

// 门店创建 同步用户关联的门店
type ShopCreatedSyncAssociatedAdmin struct {
	adminSrv *services.AdminService
}

func NewShopCreatedSyncAssociatedAdmin() *ShopCreatedSyncAssociatedAdmin {
	return &ShopCreatedSyncAssociatedAdmin{adminSrv:services.MakeAdminService()}
}

// 对应的触发事件
func (ShopCreatedSyncAssociatedAdmin) Event() message.Event {
	return events.ShopCreated{}
}

// 队列名称
func (ShopCreatedSyncAssociatedAdmin) QueueName() string {
	return "ShopCreatedSyncAssociatedAdmin"
}

// 错误处理
func (ShopCreatedSyncAssociatedAdmin) OnError(err error) {
	log.Printf("ShopCreatedSyncAssociatedAdmin Error:%s\n", err)
}

// 处理逻辑
func (this ShopCreatedSyncAssociatedAdmin) Handler(data []byte) error {
	var shop models.Shop
	err := json.Unmarshal(data, &shop)
	if err != nil {
		return err
	}
	return this.adminSrv.SyncAssociatedShop(context.Background(), &shop)
}
