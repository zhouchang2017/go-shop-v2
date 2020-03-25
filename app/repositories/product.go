package repositories

import (
	"context"
	"fmt"
	"go-shop-v2/app/models"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ProductRep struct {
	repository.IRepository
	itemRep *ItemRep
}

func ProductCacheKey(id string) string {
	return fmt.Sprintf("products:%s", id)
}

func ItemCacheKey(id string) string {
	return fmt.Sprintf("items:%s", id)
}

func NewProductRep(rep repository.IRepository, itemRep *ItemRep) *ProductRep {
	return &ProductRep{
		IRepository: rep,
		itemRep:     itemRep,
	}
}

func (this *ProductRep) GetItemRep() *ItemRep {
	return this.itemRep
}

func (this *ProductRep) FindItemById(ctx context.Context, id string) (item *models.Item, err error) {
	byId := <-this.itemRep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}

	item = byId.Result.(*models.Item)

	return item, nil
}

func (this *ProductRep) WithItems(ctx context.Context, id string) (product *models.Product, err error) {
	res := <-this.FindById(ctx, id)
	if res.Error != nil {
		return nil, res.Error
	}
	product = res.Result.(*models.Product)
	itemRes := <-this.itemRep.FindByProductId(ctx, id)
	if itemRes.Error != nil {
		return nil, itemRes.Error
	}
	product.Items = itemRes.Result.([]*models.Item)
	return product, nil
}


// 重写delete方法
func (this *ProductRep) Delete(ctx context.Context, id string) <-chan error {
	result := make(chan error)
	go func() {
		defer close(result)
		var err error

		// 开启事务
		var session mongo.Session
		if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
			result <- err
			return
		}
		if err = session.StartTransaction(); err != nil {
			result <- err
			return
		}

		err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
			objId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				session.AbortTransaction(sessionContext)
				return err
			}
			force := ctx2.GetForce(ctx)

			if force {
				// 硬删除
				_, err = this.IRepository.Collection().DeleteOne(sessionContext, bson.M{"_id": objId})
				if err != nil {
					session.AbortTransaction(sessionContext)
					return err
				}

				_, err = this.itemRep.Collection().DeleteMany(sessionContext, bson.M{
					"product.id": id,
				})
				if err != nil {
					session.AbortTransaction(sessionContext)
					return err
				}
				session.CommitTransaction(sessionContext)
			} else {
				// 软删除
				now := time.Now()
				if _, err := this.Collection().UpdateOne(sessionContext, bson.M{"_id": objId}, bson.M{
					"$set": bson.M{"deleted_at": now},
					"$currentDate": bson.M{
						"updated_at": true,
					},
				}); err != nil {
					session.AbortTransaction(sessionContext)
					return err
				}
				if _, err := this.itemRep.Collection().UpdateMany(sessionContext, bson.M{
					"product.id": id,
				}, bson.M{
					"$set": bson.M{"deleted_at": now},
					"$currentDate": bson.M{
						"updated_at": true,
					},
				}); err != nil {
					session.AbortTransaction(sessionContext)
					return err
				}
				session.CommitTransaction(sessionContext)
			}

			return nil
		})
		session.EndSession(ctx)
		result <- err
	}()
	return result
}

// 重写restore方法
func (this *ProductRep) Restore(ctx context.Context, id string) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		var err error

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			result <- repository.QueryResult{Error: err2.Err404}
			return
		}

		// 开启事务
		var session mongo.Session
		if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
			result <- repository.QueryResult{Error: err}
			return
		}
		if err = session.StartTransaction(); err != nil {
			result <- repository.QueryResult{Error: err}
			return
		}
		err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
			update := this.Collection().FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{
				"$set": bson.M{"deleted_at": nil},
				"$currentDate": bson.M{
					"updated_at": true,
				},
			})
			if update.Err() != nil {
				session.AbortTransaction(sessionContext)
				return update.Err()
			}

			if _, err = this.itemRep.Collection().UpdateMany(sessionContext, bson.M{
				"product.id": id,
			}, bson.M{
				"$set": bson.M{"deleted_at": nil},
				"$currentDate": bson.M{
					"updated_at": true,
				},
			}); err != nil {
				session.AbortTransaction(sessionContext)
				return err
			}
			result <- repository.QueryResult{Error: nil}
			session.CommitTransaction(sessionContext)
			return nil
		})
		if err != nil {
			result <- repository.QueryResult{Error: err}
		}
		session.EndSession(ctx)
	}()
	return result
}

func (this *ProductRep) AvailableOptionNames(ctx context.Context) (names []string) {
	pipline := mongo.Pipeline{
		bson.D{{"$unwind", "$options"}},
		bson.D{{"$group", bson.M{
			"_id": "$options.name",
		}}},
	}
	aggregate, err := this.Collection().Aggregate(ctx, pipline)
	if err != nil {
		return
	}
	for aggregate.Next(ctx) {
		if value, err := aggregate.Current.LookupErr("_id"); err == nil {
			name := value.StringValue()
			if name != "" {
				names = append(names, name)
			}
		}
	}
	return
}
