package repositories

import (
	"context"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const shopCartCacheKey = "shopCart"

func getShopCartCacheKey(userId string) string {
	return fmt.Sprintf("%s:%s", shopCartCacheKey, userId)
}

type ShopCartRep struct {
	*mongoRep
}


// 全选/取消
func (this *ShopCartRep) CheckedOrCancelAll(ctx context.Context, checked bool, ids ...string) (err error) {
	var objIds []primitive.ObjectID
	for _, id := range ids {
		if objId, err := primitive.ObjectIDFromHex(id); err == nil {
			objIds = append(objIds, objId)
		}
	}
	if len(objIds) > 0 {
		_, err = this.Collection().UpdateMany(ctx, bson.M{"_id": bson.M{"$in": objIds}}, bson.M{
			"$set": bson.M{"checked": checked},
			"$currentDate": bson.M{
				"updated_at": true,
			}}, )
		return err
	}
	return nil
}

func NewShopCartRep(con *mongodb.Connection) *ShopCartRep {
	return &ShopCartRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.ShopCart{}, con),
	}
}
