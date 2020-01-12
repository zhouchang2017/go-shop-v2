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

func init() {
	// 注入容器
	register(NewSyncShopAssociatedAdmin)
}

type SyncShopAssociatedAdmin struct {
	shopService *services.ShopService
}

func (SyncShopAssociatedAdmin) Event() message.Event {
	return events.AdminCreated{}
}

func (SyncShopAssociatedAdmin) QueueName() string {
	return "SyncShopAssociatedAdmin"
}

func (SyncShopAssociatedAdmin) OnError(err error) {
	log.Printf("SyncShopAssociatedAdmin error:%s\n", err)
}

func (s SyncShopAssociatedAdmin) Handler(data []byte) error {
	var admin *models.Admin
	err := json.Unmarshal(data, admin)
	if err != nil {
		return err
	}
	if admin != nil {
		return s.shopService.SyncAssociatedMembers(context.Background(), admin)
	}
	return nil
}

func NewSyncShopAssociatedAdmin() message.Listener {
	return &SyncShopAssociatedAdmin{shopService: services.MakeShopService()}
}
