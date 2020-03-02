package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type ShopCartService struct {
	rep *repositories.ShopCartRep
}

// 加入到购物车
func (this *ShopCartService) Add(ctx context.Context, userId string, itemId string, qty int64) (err error) {
	return this.rep.Add(ctx, userId, itemId, qty)
}

// 更新购物车商品itemId
// status = 1，直接更新
// status = 2, 购物车存在重复商品，进行覆盖，之前的被删除
func (this *ShopCartService) Update(ctx context.Context, userId string, beforeItemId string, afterItemId string, qty int64) (status int64, err error) {
	return this.rep.Update(ctx, userId, beforeItemId, afterItemId, qty)
}

// 购物车中物品数量增加
func (this *ShopCartService) UpdateQty(ctx context.Context, userId string, ItemId string, qty int64) (err error) {
	return this.rep.UpdateQty(ctx, userId, ItemId, qty)
}

// 全选/取消全选
func (this *ShopCartService) Toggle(ctx context.Context, userId string, checked bool, itemIds ...string) (err error) {
	return this.rep.Toggle(ctx, userId, checked, itemIds...)
}

// 删除购物车
func (this *ShopCartService) Delete(ctx context.Context, userId string, itemIds ...string) (err error) {
	return this.rep.Remove(ctx, userId, itemIds...)
}

// 前台小程序列表
func (this *ShopCartService) Index(ctx context.Context, userId string, page int64, perPage int64) (items []*models.ShopCartItem, pagination response.Pagination, err error) {
	return this.rep.Index(ctx, userId, page, perPage)
}

// 个人购物车总数
func (this *ShopCartService) Count(ctx context.Context, userId string) (count int64) {
	return this.rep.Count(ctx, userId)
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
