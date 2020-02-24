package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
)

type ShopCartService struct {
	rep *repositories.ShopCartRep
}

// 加入到购物车
func (this *ShopCartService) Add(ctx context.Context, userId string, item *models.Item, qty int64, check bool) (shopCart *models.ShopCart, err error) {
	model := models.NewShopCart(userId, item, qty, check)
	one := <-this.rep.FindOne(ctx, map[string]interface{}{
		"user_id": userId,
		"item.id": item.GetID(),
	})
	if one.Error == nil {
		cart := one.Result.(*models.ShopCart)
		updated := <-this.rep.Update(ctx, cart.GetID(), bson.M{
			"$inc": bson.M{"qty": qty},
			"$set": bson.M{
				"checked": check,
			},
		})

		if updated.Error != nil {
			return nil, updated.Error
		}
		return updated.Result.(*models.ShopCart), nil
	}
	created := <-this.rep.Create(ctx, model)
	if created.Error != nil {
		return nil, created.Error
	}
	return created.Result.(*models.ShopCart), nil
}

// 更新购物车
func (this *ShopCartService) Update(ctx context.Context, id string, qty int64, check bool) (shopCart *models.ShopCart, err error) {
	updated := <-this.rep.Update(ctx, id, bson.M{
		"$set": bson.M{
			"qty":     qty,
			"checked": check,
		},
	})
	if updated.Error != nil {
		return nil, updated.Error
	}
	return updated.Result.(*models.ShopCart), nil
}

// 删除购物车
func (this *ShopCartService) Delete(ctx context.Context, ids ...string) (err error) {
	err = <-this.rep.DeleteMany(ctx, ids...)
	return err
}

// 分页
func (this *ShopCartService) Pagination(ctx context.Context, req *request.IndexRequest) (shopCarts []*models.ShopCart, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	shopCarts = results.Result.([]*models.ShopCart)
	pagination = results.Pagination
	return
}


func NewShopCartService(rep *repositories.ShopCartRep) *ShopCartService {
	return &ShopCartService{rep: rep}
}
