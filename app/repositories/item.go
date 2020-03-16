package repositories

import (
	"context"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

type ItemRep struct {
	repository.IRepository
}

// 减库存
func (this *ItemRep) DecQty(ctx context.Context, itemId string, qty int64) error {
	objId, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		return err
	}
	result := this.Collection().FindOneAndUpdate(ctx, bson.M{
		"_id": objId,
		"qty": bson.M{"$gte": qty},
	}, bson.M{
		"$inc": bson.M{"qty": -qty},
		"$currentDate": bson.M{
			"updated_at": true,
		},
	})
	return result.Err()
}

// 增减库存
func (this *ItemRep) IncqTY(ctx context.Context, itemId string, qty int64) error {
	objId, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		return err
	}
	result := this.Collection().FindOneAndUpdate(ctx, bson.M{
		"_id": objId,
	}, bson.M{
		"$inc": bson.M{"qty": qty},
		"$currentDate": bson.M{
			"updated_at": true,
		},
	})
	return result.Err()
}

// cache key : products:id items
// cache key: items:item_id detail
// cache key: items:item_id qty

func (this *ItemRep) FindByProductId(ctx context.Context, id string) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	go func() {
		defer close(output)

		//if redis.GetConFn() != nil {
		//	result, err := redis.GetConFn().HMGet(ProductCacheKey(id), "items").Result()
		//	if err == nil {
		//		if result[0] != nil {
		//			var items []*models.Item
		//			jsonValue := result[0].(string)
		//			if err := json.Unmarshal([]byte(jsonValue), &items); err != nil {
		//				log.Printf("%s [%s] FindByProductId form cache,error:%s\n", this.TableName(), id, err)
		//				// 从缓存中移除
		//				redis.GetConFn().HDel(ProductCacheKey(id), "items")
		//			}
		//			output <- repository.QueryResult{
		//				Result: items,
		//				Error:  nil,
		//			}
		//			return
		//		}
		//	}
		//}

		many := <-this.FindMany(ctx, map[string]interface{}{"product.id": id})

		//if many.Error == nil {
		//	if redis.GetConFn() != nil && many.Result != nil {
		//		items := many.Result.([]*models.Item)
		//		if marshal, err := json.Marshal(items); err == nil {
		//			redis.GetConFn().HMSet(ProductCacheKey(id), "items", marshal)
		//		}
		//		for _,item:=range items {
		//			// 库存缓存
		//			redis.GetConFn().HMSet(ItemCacheKey(item.GetID()), "qty", item.Qty)
		//		}
		//	}
		//}
		output <- repository.QueryResult{
			Result: many.Result,
			Error:  many.Error,
		}
	}()
	return output
}

func (this *ItemRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "code", Value: bsonx.Int64(-1)}},
			Options: options.Index().SetBackground(true),
		},
	}
}

func NewItemRep(rep repository.IRepository) *ItemRep {
	repository := &ItemRep{rep}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Printf("model [%s] create indexs error:%s\n", repository.TableName(), err)
		panic(err)
	}
	return repository
}
