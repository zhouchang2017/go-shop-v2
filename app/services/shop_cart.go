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
		// 已存在
		cart := one.Result.(*models.ShopCart)
		updated := <-this.rep.Update(ctx, cart.GetID(), bson.M{
			"$inc": bson.M{"qty": qty},
			"$set": bson.M{
				"checked": check,
				"price":   item.Price,
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
func (this *ShopCartService) Update(ctx context.Context, userId string, id string, item *models.Item, qty int64, check bool) (shopCart *models.ShopCart, err error) {
	if item != nil {
		// 查询是否有相同sku 购物车
		one := <-this.rep.FindOne(ctx, bson.M{
			"user_id": userId,
			"item.id": item.GetID(),
		})
		if one.Error == nil {
			// 存在相同sku，进行覆盖
			// 删除当前id
			if err := <-this.rep.Delete(ctx, id); err != nil {
				return nil, err
			}

			//shopCart = &models.ShopCart{}
			cart := one.Result.(*models.ShopCart)
			cart.Qty = qty
			cart.Checked = check
			cart.Item = item.ToAssociated()
			saved := <-this.rep.Save(ctx, cart)
			if saved.Error != nil {
				err = saved.Error
				return
			}
			return saved.Result.(*models.ShopCart), nil
		}
	}

	updated := bson.M{
		"qty":     qty,
		"checked": check,
	}
	if item != nil {
		updated["item"] = item.ToAssociated()
	}

	// 不存在相同sku
	// 直接更新
	result := <-this.rep.Update(ctx, id, bson.M{
		"$set": updated,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	return result.Result.(*models.ShopCart), nil
}

// 全选/取消全选
func (this *ShopCartService) CheckedOrCancelAll(ctx context.Context, checked bool, ids ...string) (err error) {
	return this.rep.CheckedOrCancelAll(ctx, checked, ids...)
}

// 删除购物车
func (this *ShopCartService) Delete(ctx context.Context, ids ...string) (err error) {
	err = <-this.rep.DeleteMany(ctx, ids...)
	return err
}

// 分页
func (this *ShopCartService) Pagination(ctx context.Context, req *request.IndexRequest) (shopCarts []*models.ShopCart, pagination response.Pagination, err error) {
	// 不进行分页
	req.Page = -1
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	shopCarts = results.Result.([]*models.ShopCart)
	if len(shopCarts) == 0 {
		shopCarts = []*models.ShopCart{}
	}
	pagination = results.Pagination
	return
}

func NewShopCartService(rep *repositories.ShopCartRep) *ShopCartService {
	return &ShopCartService{rep: rep}
}
