package repositories

import (
	"context"
	"go-shop-v2/app/models"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShopRep struct {
	*mongoRep
}

func (s *ShopRep) GetAllAssociatedShops(ctx context.Context) (res []*models.AssociatedShop) {
	trashed := ctx2.GetTrashed(ctx)

	filter := bson.M{}
	if !trashed {
		filter["deleted_at"] = bson.D{{"$eq", nil}}
	}

	cursor, err := s.Collection().Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 1, "name": 1}))
	if err != nil {
		return res
	}
	var shops []*models.Shop
	err = cursor.All(ctx, &shops)
	if err != nil {
		return res
	}
	for _, shop := range shops {
		res = append(res, shop.ToAssociated())
	}
	return res
}

func NewShopRep(con *mongodb.Connection) *ShopRep {
	return &ShopRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Shop{}, con),
	}
}
