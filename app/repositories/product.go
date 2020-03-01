package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sync"
)

type ProductRep struct {
	repository.IRepository
	itemRep *ItemRep
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
	res := <-this.itemRep.FindById(ctx, id)
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
		p := entity.(*models.Product)
		items := p.Items
		res := <-this.IRepository.Create(ctx, entity)
		if res.Error != nil {
			result <- repository.InsertResult{Error: res.Error}
			return
		}
		product := res.Result.(*models.Product)
		wg := sync.WaitGroup{}
		//var newItems []*models.Item
		newItems := make([]*models.Item, len(items))
		for index, item := range items {
			item.ID = primitive.NewObjectID()
			wg.Add(1)
			item.Product = product.ToAssociated()
			go func(ind int, i *models.Item) {
				defer wg.Done()
				itemRes := <-this.itemRep.Create(ctx, i)
				if itemRes.Error != nil {
					log.Printf("create item error:%s", itemRes.Error)
					result <- repository.InsertResult{Error: res.Error}
					return
				}
				newItems[ind] = itemRes.Result.(*models.Item)
			}(index, item)
		}
		wg.Wait()
		product.Items = newItems
		result <- repository.InsertResult{Id: res.Id, Result: product}
	}()
	return result
}

// 重写save方法
func (this *ProductRep) Save(ctx context.Context, entity interface{}) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		p := entity.(*models.Product)
		items := p.Items
		// 储存product
		productSaved := <-this.IRepository.Save(ctx, entity)
		if productSaved.Error != nil {
			result <- repository.QueryResult{Error: productSaved.Error}
			return
		}
		product := productSaved.Result.(*models.Product)
		// 变体更新
		wg := sync.WaitGroup{}
		var newItems []*models.Item
		for _, item := range items {
			wg.Add(1)
			item.Product = product.ToAssociated()
			go func(i *models.Item) {
				defer wg.Done()

				if i.ID.IsZero() {
					// 新增变体
					created := <-this.itemRep.Create(ctx, i)
					if created.Error != nil {
						log.Printf("update product %s add item error:%s", product.GetID(), created.Error)
						result <- repository.QueryResult{Error: created.Error}
						return
					}
					newItems = append(newItems, created.Result.(*models.Item))
				} else {
					saved := <-this.itemRep.Save(ctx, i)
					if saved.Error != nil {
						log.Printf("update product %s save item error:%s", product.GetID(), saved.Error)
						result <- repository.QueryResult{Error: saved.Error}
						return
					}
					newItems = append(newItems, saved.Result.(*models.Item))
				}

			}(item)
		}
		wg.Wait()
		product.Items = newItems
		result <- repository.QueryResult{Result: product}
	}()
	return result
}
