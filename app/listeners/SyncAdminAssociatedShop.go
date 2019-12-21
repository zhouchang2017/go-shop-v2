package listeners

import (
	"context"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
)

func init() {
	// 注入容器
	register(NewSyncAdminAssociatedShop)
}

// 同步用户关联门店
type SyncAdminAssociatedShop struct {
	adminService *services.AdminService
}

func NewSyncAdminAssociatedShop(adminService *services.AdminService) *SyncAdminAssociatedShop {
	return &SyncAdminAssociatedShop{adminService: adminService}
}

func (s SyncAdminAssociatedShop) Handle(ctx context.Context, event interface{}) error {
	var shop *models.Shop
	switch event.(type) {
	case events.ShopCreated:
		shop = event.(events.ShopCreated).Shop
	case events.ShopUpdated:
		shop = event.(events.ShopUpdated).Shop
	default:
		shop = nil
	}
	if shop != nil {
		return s.adminService.SyncAssociatedShop(ctx, shop)
	}
	return nil
}
