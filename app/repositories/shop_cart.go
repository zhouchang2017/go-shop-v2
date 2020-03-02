package repositories

import (
	"context"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const shopCartCacheKey = "shopCart"

func getShopCartCacheKey(userId string) string {
	return fmt.Sprintf("%s:%s", shopCartCacheKey, userId)
}

type ShopCartRep struct {
	repository.IRepository
}

// 分页
func (this *ShopCartRep) Index(ctx context.Context, userId string, page int64, perPage int64) (items []*models.ShopCartItem, pagination response.Pagination, err error) {
	items = []*models.ShopCartItem{}
	count := this.Count(ctx, userId)
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = 15
	}
	skip := (page - 1) * 15

	pagination = response.Pagination{
		Total:       count,
		CurrentPage: page,
		PerPage:     perPage,
		HasNextPage: page*perPage < count,
	}

	result := this.Collection().FindOne(ctx, bson.M{"user_id": userId}, options.FindOne().SetProjection(bson.M{
		"items": bson.M{"$slice": bson.A{skip, perPage * page}},
	}))

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return
		}
		err = result.Err()
		return
	}
	var res models.ShopCart
	if err := result.Decode(&res); err != nil {
		return items, pagination, err
	}
	for _, item := range res.Items {
		items = append(items, item)
	}
	return items, pagination, nil
}

// 添加数量
func (this *ShopCartRep) UpdateQty(ctx context.Context, userId string, itemId string, qty int64) (err error) {
	arrayFilters := []interface{}{
		bson.M{
			"elem.item_id": itemId,
		},
	}
	_, err = this.Collection().UpdateOne(ctx,
		bson.M{
			"user_id": userId,
		},
		bson.M{
			"$set": bson.M{"items.$[elem].qty": qty},
			"$currentDate": bson.M{
				"updated_at": true,
			},
		},
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: arrayFilters,
		}),
	)

	return
}

// 购物车商品更新状态
// updated = true 说明更新
// deleted = true 说明之前的被删除

// 更新itemId
func (this *ShopCartRep) Update(ctx context.Context, userId string, beforeItemId string, afterItemId string, qty int64) (status int64, err error) {
	result := this.Collection().FindOne(ctx,
		bson.M{
			"user_id": userId,
			"items":   bson.M{"$elemMatch": bson.M{"item_id": afterItemId}},
		},
	)
	if result.Err() == nil {
		// 当前购物车中已存在相同商品，进行覆盖

		// 删除 beforeItemId
		if err := this.Remove(ctx, userId, beforeItemId); err != nil {
			return status, err
		}
		if err := this.UpdateQty(ctx, userId, afterItemId, qty); err != nil {
			return status, err
		}
		status = 2
		return
	}

	// 购物车中不存在，需要更新
	arrayFilters := []interface{}{
		bson.M{
			"elem.item_id": beforeItemId,
		},
	}
	_, err = this.Collection().UpdateOne(ctx,
		bson.M{
			"user_id": userId,
		},
		bson.M{
			"$set": bson.M{"items.$[elem].item_id": afterItemId, "items.$[elem].qty": qty},
			"$currentDate": bson.M{
				"updated_at": true,
			},
		},
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: arrayFilters,
		}),
	)

	status = 1

	return
}

// 添加
func (this *ShopCartRep) Add(ctx context.Context, userId string, itemId string, qty int64) (err error) {

	result := this.Collection().FindOne(ctx,
		bson.M{
			"user_id": userId,
			"items":   bson.M{"$elemMatch": bson.M{"item_id": itemId}},
		},
	)
	if result.Err() == nil {
		_, err = this.Collection().UpdateOne(ctx,
			bson.M{
				"user_id": userId,
				"items":   bson.M{"$elemMatch": bson.M{"item_id": itemId}},
			},
			bson.M{
				"$inc": bson.M{
					"items.$.qty": qty,
				},
				"$currentDate": bson.M{
					"updated_at": true,
				},
			})
		return err
	}

	// 存在错误，数据库不存在
	// 新增

	_, err = this.Collection().UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{
		"$push": bson.M{
			"items": bson.M{
				"$each":     bson.A{bson.M{"item_id": itemId, "qty": qty, "checked": true}},
				"$position": 0,
			},
		},
		"$currentDate": bson.M{
			"updated_at": true,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}, options.Update().SetUpsert(true))
	return
}

// 移除
func (this *ShopCartRep) Remove(ctx context.Context, userId string, itemIds ...string) (err error) {
	_, err = this.Collection().UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{
		"$pull": bson.M{
			"items": bson.M{"item_id": bson.M{"$in": itemIds}},
		},
		"$currentDate": bson.M{
			"updated_at": true,
		},
	})
	return err
}

// 全选/取消
func (this *ShopCartRep) Toggle(ctx context.Context, userId string, checked bool, itemIds ...string) (err error) {
	arrayFilters := []interface{}{
		bson.M{
			"elem.item_id": bson.M{"$in": itemIds},
		},
	}
	_, err = this.Collection().UpdateOne(ctx,
		bson.M{
			"user_id": userId,
		},
		bson.M{
			"$set": bson.M{"items.$[elem].checked": checked},
			"$currentDate": bson.M{
				"updated_at": true,
			},
		},
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: arrayFilters,
		}),
	)

	return
}

// 个人购物车商品总数
func (this *ShopCartRep) Count(ctx context.Context, userId string) (count int64) {
	aggregate, err := this.Collection().Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match", bson.M{"user_id": userId}}},
		bson.D{{"$project", bson.M{"count": bson.M{"$size": "$items"}}}},
	})
	if err != nil {
		return 0
	}
	var res []countRep
	err = aggregate.All(ctx, &res)
	if err != nil {
		return 0
	}
	if len(res) > 0 {
		return res[0].Count
	}
	return 0
}

func NewShopCartRep(rep repository.IRepository) *ShopCartRep {
	return &ShopCartRep{rep}
}
