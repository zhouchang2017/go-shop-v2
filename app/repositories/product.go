package repositories

import (
	"context"
	"fmt"
	"go-shop-v2/app/models"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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

// 重写Create方法
func (this *ProductRep) Create(ctx context.Context, entity interface{}) <-chan repository.InsertResult {
	result := make(chan repository.InsertResult)
	go func() {
		defer close(result)

		// 开启事务
		var session mongo.Session
		var err error
		if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
			result <- repository.InsertResult{Error: err}
			return
		}
		if err = session.StartTransaction(); err != nil {
			result <- repository.InsertResult{Error: err}
			return
		}

		err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
			p := entity.(*models.Product)
			p.SetAvatar()
			items := p.Items
			res := <-this.IRepository.Create(sessionContext, entity)
			if res.Error != nil {
				session.AbortTransaction(sessionContext)
				return res.Error
			}
			product := res.Result.(*models.Product)
			newItems := []*models.Item{}
			for _, item := range items {
				item.Product = product.ToAssociated()
				item.SetAvatar()
				itemRes := <-this.itemRep.Create(sessionContext, item)
				if itemRes.Error != nil {
					log.Printf("create item error:%s", itemRes.Error)
					session.AbortTransaction(sessionContext)
					return itemRes.Error
				}
				newItems = append(newItems, itemRes.Result.(*models.Item))
			}
			product.Items = newItems
			result <- repository.InsertResult{Id: res.Id, Result: product}
			session.CommitTransaction(sessionContext)
			return nil
		})

		session.EndSession(ctx)
		if err != nil {
			result <- repository.InsertResult{
				Error: err,
			}
		}

	}()
	return result
}

// 重写save方法
func (this *ProductRep) Save(ctx context.Context, entity interface{}) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)

		// 开启事务
		var session mongo.Session
		var err error
		if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
			result <- repository.QueryResult{Error: err}
			return
		}
		if err = session.StartTransaction(); err != nil {
			result <- repository.QueryResult{Error: err}
			return
		}
		p := entity.(*models.Product)
		p.SetAvatar()
		err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {

			items := p.Items
			// 储存product
			productSaved := <-this.IRepository.Save(sessionContext, entity)
			if productSaved.Error != nil {
				session.AbortTransaction(sessionContext)
				return productSaved.Error
			}
			product := productSaved.Result.(*models.Product)
			// 变体更新
			var newItems []*models.Item
			for _, item := range items {
				item.Product = product.ToAssociated()
				item.SetAvatar()
				if item.ID.IsZero() {
					// 新增变体
					created := <-this.itemRep.Create(sessionContext, item)
					if created.Error != nil {
						log.Printf("update product %s add item error:%s", product.GetID(), created.Error)
						session.AbortTransaction(sessionContext)
						return created.Error
					}
					newItems = append(newItems, created.Result.(*models.Item))
				} else {
					saved := <-this.itemRep.Save(sessionContext, item)
					if saved.Error != nil {
						log.Printf("update product %s save item error:%s", product.GetID(), saved.Error)
						session.AbortTransaction(sessionContext)
						return saved.Error
					}
					newItems = append(newItems, saved.Result.(*models.Item))
				}
			}
			product.Items = newItems
			result <- repository.QueryResult{Result: product}
			session.CommitTransaction(sessionContext)
			return nil
		})
		session.EndSession(ctx)
		if err != nil {
			result <- repository.QueryResult{Error: err}
		} else {
			// 更新缓存
			//if redis.GetConFn() != nil {
			//	if p != nil {
			//		// product detail缓存
			//		if marshal, err := json.Marshal(p); err == nil {
			//			redis.GetConFn().HMSet(ProductCacheKey(p.GetID()), "detail", marshal)
			//		}
			//
			//		// items缓存
			//		if marshal, err := json.Marshal(p.Items); err == nil {
			//			redis.GetConFn().HMSet(ProductCacheKey(p.GetID()), "items", marshal)
			//		}
			//
			//		for _, item := range p.Items {
			//			// 库存缓存
			//			redis.GetConFn().HMSet(ItemCacheKey(item.GetID()), "qty", item.Qty)
			//		}
			//	}
			//}
		}
	}()
	return result
}

// 重写delete方法
func (this *ProductRep) Delete(ctx context.Context, id string) <-chan error {
	result := make(chan error)
	go func() {
		defer close(result)

		//var product *models.Product
		var err error
		//if redis.GetConFn() != nil {
		//	product, err = this.WithItems(ctx, id)
		//	if err != nil {
		//		result <- err
		//		return
		//	}
		//
		//}

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
		//if err == nil {
		//	// 清除缓存
		//	if redis.GetConFn() != nil {
		//		if product != nil {
		//			redis.GetConFn().HDel(ProductCacheKey(product.GetID()), "detail")
		//			redis.GetConFn().HDel(ProductCacheKey(product.GetID()), "items")
		//			// 库存缓存
		//			for _, item := range product.Items {
		//				redis.GetConFn().HDel(ItemCacheKey(item.GetID()), "qty")
		//			}
		//		}
		//	}
		//}
		result <- err
	}()
	return result
}
