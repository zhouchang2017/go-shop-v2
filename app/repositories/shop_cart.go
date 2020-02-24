package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

type ShopCartRep struct {
	*mongoRep
}

// 从购物车中删除
//func (this *ShopCartRep) RemoveItems(ctx context.Context, ids ...string) (err error) {
//
//}

func NewShopCartRep(con *mongodb.Connection) *ShopCartRep {
	return &ShopCartRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.ShopCart{}, con),
	}
}
