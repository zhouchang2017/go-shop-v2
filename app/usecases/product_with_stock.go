package usecases

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go.mongodb.org/mongo-driver/bson"
)


// 产品详情
func ProductWithStock(ctx context.Context, id string, productSrv *services.ProductService, inventorySrv *services.InventoryService) (product *models.Product, err error) {
	product, err = productSrv.FindByIdWithItems(ctx, id)
	if err != nil {
		return
	}

	inventories, err := inventorySrv.Search(ctx, &repositories.QueryOption{
		ProductId: id,
		Qty:       bson.M{"$gt": 0},
		Status:    []int8{0},
	})
	if err == nil {
		for _, inventory := range inventories {
			for _, item := range product.Items {
				if inventory.Item.Id == item.GetID() {
					product.Qty += inventory.Qty
					item.Qty += inventory.Qty
					item.Inventories = append(item.Inventories, inventory)
				}
			}
		}
	}



	return
}
