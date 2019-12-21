package listeners

import (
	"context"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
)

func init() {
	// 注入容器
	register(NewSyncShopAssociatedAdmin)
}

type SyncShopAssociatedAdmin struct {
	shopService *services.ShopService
}

func NewSyncShopAssociatedAdmin(shopService *services.ShopService) *SyncShopAssociatedAdmin {
	return &SyncShopAssociatedAdmin{shopService: shopService}
}

func (s SyncShopAssociatedAdmin) Handle(ctx context.Context, event interface{}) error {
	var admin *models.Admin
	switch event.(type) {
	case events.AdminCreated:
		admin = event.(events.AdminCreated).Admin
	case events.AdminUpdated:
		admin = event.(events.AdminUpdated).Admin
	default:
		admin = nil
	}
	if admin != nil {
		return s.shopService.SyncAssociatedMembers(ctx, admin)
	}
	return nil
}
