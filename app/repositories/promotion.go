package repositories

import (
	"context"
	"fmt"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"time"
)

type PromotionRep struct {
	repository.IRepository
	promotionItemRep *PromotionItemRep
}


func (this *PromotionRep) FindByIdWithItems(ctx context.Context, id string) (promotion *models.Promotion, err error) {

	results := <-this.FindById(ctx, id)
	if results.Error != nil {
		return nil, results.Error
	}
	promotion = results.Result.(*models.Promotion)

	find, err := this.promotionItemRep.Collection().Find(ctx, bson.M{"promotion.id": id})
	if err != nil {
		return nil, err
	}
	items := []*models.PromotionItem{}
	if err := find.All(ctx, &items); err != nil {
		return nil, err
	}
	promotion.Items = items
	return promotion, nil
}

func productPriceCacheKey(id string) string {
	return fmt.Sprintf("products:%s", id)
}

// 获取Product显示价格
// 显示sku最低价格
// 如果产品,不存在促进价格 返回-1
// 缓存key products:id price
func (this *PromotionRep) FindProductsPrice(ctx context.Context, productIds ...string) (prices map[string]int64) {

	var group errgroup.Group
	prices = map[string]int64{}
	sem := make(chan struct{}, 10)
	for _, id := range productIds {
		//prices[id] = -1
		id := id

		if redis.GetConFn() != nil {
			result, err := redis.GetConFn().HMGet(productPriceCacheKey(id), "price").Result()
			if err == nil {
				if len(result) > 0 {
					price := result[0]
					if price != nil {
						// hit
						prices[id] = price.(int64)
						break
					}

				}
			}
		}

		sem <- struct{}{}
		group.Go(func() error {
			promotionItem := &models.PromotionItem{}
			result := this.promotionItemRep.Collection().FindOne(ctx, bson.M{
				"product_id":         id,
				"promotion.type":     models.UnitSale, // 单品优惠
				"promotion.enable":   true,
				"promotion.begin_at": bson.M{"$lte": time.Now()},
				"promotion.ended_at": bson.M{"$gt": time.Now()},
			})
			if result.Err() != nil {
				// 不存在促销
				// set cache
				if redis.GetConFn() != nil {
					redis.GetConFn().HMSet(productPriceCacheKey(id), "price", -1)
				}
				prices[id] = -1
			}
			err := result.Decode(promotionItem)
			if err == nil {
				var p []int64
				for _, item := range promotionItem.Units {
					p = append(p, item.Price)
				}
				minPrice := utils.Min(p...)
				prices[id] = minPrice
				if redis.GetConFn() != nil {
					redis.GetConFn().HMSet(productPriceCacheKey(id), "price", minPrice)

				}
			}
			<-sem
			return err
		})
	}
	if err := group.Wait(); err != nil {
		return prices
	}
	return prices

}

// 获取Product对应的活动
func (this *PromotionRep) FindActivePromotionByProductId(ctx context.Context, productId string) (items []*models.PromotionItem) {
	items = []*models.PromotionItem{}
	find, err := this.promotionItemRep.Collection().Find(ctx, bson.M{
		"product_id":         productId,
		"promotion.enable":   true,
		"promotion.begin_at": bson.M{"$lte": time.Now()},
		"promotion.ended_at": bson.M{"$gt": time.Now()},
	})
	if err != nil {
		return items
	}
	if err := find.All(ctx, &items); err != nil {
		return items
	}
	return items
}

// 重写Create
func (this *PromotionRep) Create(ctx context.Context, entity interface{}) <-chan repository.InsertResult {
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
			p := entity.(*models.Promotion)
			items := p.Items

			res := <-this.IRepository.Create(sessionContext, entity)
			if res.Error != nil {
				result <- repository.InsertResult{Error: res.Error}
				session.AbortTransaction(sessionContext)
				return res.Error
			}
			// 创建items
			promotion := res.Result.(*models.Promotion)
			promotionItems := make([]*models.PromotionItem, len(items))

			for index, item := range items {
				item.Promotion = promotion.ToAssociated()
				created := <-this.promotionItemRep.Create(sessionContext, item)
				if created.Error != nil {
					result <- repository.InsertResult{Error: created.Error}
					session.AbortTransaction(sessionContext)
					return created.Error
				}
				promotionItems[index] = created.Result.(*models.PromotionItem)
			}

			promotion.Items = promotionItems
			result <- repository.InsertResult{
				Id:     promotion.GetID(),
				Result: promotion,
				Error:  nil,
			}
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

// 获取Promotion -> Items -> []ids
func (this *PromotionRep) resolveItemIds(ctx context.Context, id string) (ids []string) {
	ids = []string{}
	find, err := this.promotionItemRep.Collection().Find(ctx, bson.M{"promotion.id": id}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return
	}
	defer find.Close(ctx)
	for find.Next(ctx) {
		lookup := find.Current.Lookup("_id")
		var id primitive.ObjectID
		err := lookup.Unmarshal(&id)
		if err != nil {
			return
		}
		ids = append(ids, id.Hex())
	}
	return
}

// 重写Save
func (this *PromotionRep) Save(ctx context.Context, entity interface{}) <-chan repository.QueryResult {
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

		err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
			p := entity.(*models.Promotion)
			items := p.Items
			// promotion
			promotionSaved := <-this.IRepository.Save(ctx, entity)
			if promotionSaved.Error != nil {
				session.AbortTransaction(sessionContext)
				return promotionSaved.Error
			}
			promotion := promotionSaved.Result.(*models.Promotion)

			// 查询所有items ids
			ids := this.resolveItemIds(ctx, promotion.GetID())

			var deleteIds []string
			promotionItems := []*models.PromotionItem{}
			for _, item := range items {
				if item.ID.IsZero() {
					// 新增
					item.Promotion = promotion.ToAssociated()
					created := <-this.promotionItemRep.Create(ctx, item)
					if created.Error != nil {
						session.AbortTransaction(sessionContext)
						return created.Error
					}
					promotionItems = append(promotionItems, created.Result.(*models.PromotionItem))
				} else {
					var exist = false
					for _, id := range ids {
						if id == item.GetID() {
							exist = true
						}
					}
					if !exist {
						deleteIds = append(deleteIds, item.GetID())
						break
					}
					// 更新
					item.Promotion = promotion.ToAssociated()
					saved := <-this.promotionItemRep.Save(ctx, item)
					if saved.Error != nil {
						session.AbortTransaction(sessionContext)
						return saved.Error
					}
					promotionItems = append(promotionItems, saved.Result.(*models.PromotionItem))
				}
			}

			if err = <-this.promotionItemRep.DeleteMany(ctx, deleteIds...); err != nil {
				// 删除失败
				session.AbortTransaction(sessionContext)
				return err
			}

			promotion.Items = promotionItems
			result <- repository.QueryResult{
				Result: promotion,
			}
			session.CommitTransaction(sessionContext)
			return nil
		})
		session.EndSession(ctx)
		if err != nil {
			result <- repository.QueryResult{Error: err}
		}
	}()
	return result
}

// 重新Delete
func (this *PromotionRep) Delete(ctx context.Context, id string) <-chan error {
	result := make(chan error)
	go func() {
		defer close(result)

		// 开启事务
		var session mongo.Session
		var err error
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
			_, err = this.IRepository.Collection().DeleteOne(sessionContext, bson.M{"_id": objId})
			if err != nil {
				session.AbortTransaction(sessionContext)
				return err
			}

			_, err = this.promotionItemRep.Collection().DeleteMany(sessionContext, bson.M{
				"promotion.id": id,
			})
			if err != nil {
				session.AbortTransaction(sessionContext)
				return err
			}
			session.CommitTransaction(sessionContext)
			return nil
		})
		session.EndSession(ctx)
		result <- err
	}()
	return result
	//
}

func NewPromotionRep(IRepository repository.IRepository, promotionItemRep *PromotionItemRep) *PromotionRep {
	return &PromotionRep{IRepository: IRepository, promotionItemRep: promotionItemRep}
}