package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

func init() {
	register(NewShopRep)
}

type ShopRep struct {
	*mongoRep
}

// 添加成员
func (s *ShopRep) AddMember(ctx context.Context) {

}

// 删除成员
func (s *ShopRep) DeleteMember(ctx context.Context) {

}

// 更新成员

func NewShopRep(con *mongodb.Connection) *ShopRep {
	return &ShopRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Shop{}, con),
	}
}
