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
	register(NewSyncAdminAssociatedShop)
}

// 同步用户关联门店
type SyncAdminAssociatedShop struct {
	adminService *services.AdminService
}

func (SyncAdminAssociatedShop) Event() message.Event {
	return events.ShopCreated{}
}

func (SyncAdminAssociatedShop) QueueName() string {
	return "SyncAdminAssociatedShop"
}

func (SyncAdminAssociatedShop) OnError(err error) {
	log.Printf("SyncAdminAssociatedShop error:%s\n", err)
}

func (s SyncAdminAssociatedShop) Handler(data []byte) error {
	var shop *models.Shop
	err := json.Unmarshal(data, shop)
	if err != nil {
		return err
	}
	if shop != nil {
		return s.adminService.SyncAssociatedShop(context.Background(), shop)
	}
	return nil
}

func NewSyncAdminAssociatedShop() message.Listener {
	return &SyncAdminAssociatedShop{adminService: services.MakeAdminService()}
}

//func (s SyncAdminAssociatedShop) Handle(ctx context.Context, event interface{}) error {
//	var shop *models.Shop
//	switch event.(type) {
//	case events.ShopCreated:
//		shop = event.(events.ShopCreated).Shop
//	case events.ShopUpdated:
//		shop = event.(events.ShopUpdated).Shop
//	default:
//		shop = nil
//	}
//	if shop != nil {
//		return s.adminService.SyncAssociatedShop(ctx, shop)
//	}
//	return nil
//}
