package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type ShopCartService struct {
	rep          *repositories.ShopCartRep
	itemRep      *repositories.ItemRep
	promotionRep *repositories.PromotionRep
}

// 加入到购物车
func (this *ShopCartService) Add(ctx context.Context, userId string, itemId string, qty int64) (shopCartItem *models.ShopCartItem, err error) {
	item, promotions, err := this.findItem(ctx, itemId)
	if err != nil {
		return
	}

	if item.Qty < qty {
		// 数量不足
		return nil, err2.Err422.F("库存不足")
	}

	err = this.rep.Add(ctx, userId, item, qty)

	return &models.ShopCartItem{
		ItemId:     itemId,
		Item:       item,
		Promotions: promotions,
		Price:      item.PromotionPrice,
		Qty:        qty,
	}, err
}

func (this *ShopCartService) findItem(ctx context.Context, id string) (item *models.Item, promotions []*models.PromotionItem, err error) {
	result := <-this.itemRep.FindById(ctx, id)
	if result.Error != nil {
		// 产品不存在
		err = result.Error
		return
	}
	item = result.Result.(*models.Item)
	item.PromotionPrice = item.Price
	// 获取产品促销活动
	promotions = this.promotionRep.FindActivePromotionByProductId(ctx, item.Product.Id)
	if len(promotions) > 0 {
		if promotions[0].Promotion.Type == models.UnitSale {
			price := promotions[0].FindPriceByItemId(item.GetID())
			if price != -1 {
				item.PromotionPrice = price
			}
		}
	}
	return item, promotions, nil
}

// 更新购物车商品itemId
// status = 1，直接更新
// status = 2, 购物车存在重复商品，进行覆盖，之前的被删除
func (this *ShopCartService) Update(ctx context.Context, userId string, beforeItemId string, afterItemId string, qty int64) (status int64, shopCartItem *models.ShopCartItem, err error) {
	afterItem, promotions, err := this.findItem(ctx, afterItemId)
	if err != nil {
		return
	}
	status, err = this.rep.Update(ctx, userId, beforeItemId, afterItem, qty)
	if err != nil {
		return
	}
	shopCartItem = &models.ShopCartItem{
		ItemId:     afterItemId,
		Item:       afterItem,
		Promotions: promotions,
		Price:      afterItem.PromotionPrice,
		Qty:        qty,
	}
	return
}

// 购物车中物品数量增加
func (this *ShopCartService) UpdateQty(ctx context.Context, userId string, itemId string, qty int64) (shopCartItem *models.ShopCartItem, err error) {
	item, promotions, err := this.findItem(ctx, itemId)
	if err != nil {
		return
	}
	if item.Qty < qty {
		err = err2.Err422.F("库存不足")
		return
	}
	err = this.rep.UpdateQty(ctx, userId, item, qty)
	if err != nil {
		return
	}
	return &models.ShopCartItem{
		ItemId:     itemId,
		Item:       item,
		Promotions: promotions,
		Price:      item.PromotionPrice,
		Qty:        qty,
	}, nil
}

// 全选/取消全选
func (this *ShopCartService) Toggle(ctx context.Context, userId string, checked bool, itemIds ...string) (err error) {
	return this.rep.Toggle(ctx, userId, checked, itemIds...)
}

// 删除购物车
func (this *ShopCartService) Delete(ctx context.Context, userId string, itemIds ...string) (err error) {
	return this.rep.Remove(ctx, userId, itemIds...)
}

// 通过itemId匹配item实体
func (this ShopCartService) resolveItem(items []*models.Item, id string) *models.Item {
	for _, item := range items {
		if item.GetID() == id {
			return item
		}
	}
	return nil
}

// 前台小程序列表
func (this *ShopCartService) Index(ctx context.Context, userId string, page int64, perPage int64) (items []*models.ShopCartItem, pagination response.Pagination, err error) {
	items, pagination, err = this.rep.Index(ctx, userId, page, perPage)
	if err != nil {
		return
	}
	// 加载items
	var itemIds []string
	for _, item := range items {
		itemIds = append(itemIds, item.ItemId)
	}
	// 加载软删除变体，前端标记为已失效
	withTrashedCtx := ctx2.WithTrashed(ctx, true)
	result := <-this.itemRep.FindByIds(withTrashedCtx, itemIds...)
	if result.Error != nil {
		err = result.Error
		return
	}

	// productIdsMap := map[string][]*models.ShopCartItem{}

	itemEntities := result.Result.([]*models.Item)
	for _, item := range items {
		// 设置默认促销活动空数组
		item.Promotions = []*models.PromotionItem{}
		if resolveItem := this.resolveItem(itemEntities, item.ItemId); resolveItem != nil {
			item.Item = resolveItem
			// 设置默认促销价格
			item.Item.PromotionPrice = item.Item.Price
			//if _, ok := productIdsMap[item.Item.Product.Id]; ok {
			//	productIdsMap[item.Item.Product.Id] = append(productIdsMap[item.Item.Product.Id], item)
			//} else {
			//	productIdsMap[item.Item.Product.Id] = []*models.ShopCartItem{item}
			//}
		}
	}

	// 加载促销信息
	//for productId, items := range productIdsMap {
	//	promotionItems := this.promotionRep.FindActivePromotionByProductId(ctx, productId)
	//
	//	for _, item := range items {
	//		// 判断item是否属于活动
	//		item.Promotions = models.PromotionItems(promotionItems).IncludeItem(item.ItemId, productId)
	//
	//		if len(promotionItems) > 0 {
	//			if promotionItems[0].Promotion.Type == models.UnitSale {
	//				price := promotionItems[0].FindPriceByItemId(item.ItemId)
	//				if price != -1 {
	//					item.Item.PromotionPrice = price
	//				}
	//			}
	//		}
	//	}
	//}

	return
}

// 个人购物车总数
func (this *ShopCartService) Count(ctx context.Context, userId string) (count int64) {
	return this.rep.Count(ctx, userId)
}

type shopCartItemsDetail struct {
	Item       *models.Item            `json:"item"`
	Promotions []*models.PromotionItem `json:"promotions"`
}

// 获取checked商品的详细信息
func (this *ShopCartService) GetShopCartItemsDetail(ctx context.Context, userId string, itemIds ...string) (details map[string]*shopCartItemsDetail, err error) {
	result := <-this.itemRep.FindByIds(ctx, itemIds...)
	if result.Error != nil {
		err = result.Error
		return
	}

	productIdsMap := map[string][]*shopCartItemsDetail{}

	items := result.Result.([]*models.Item)
	for _, item := range items {
		item.PromotionPrice = item.Price
		if _, ok := productIdsMap[item.Product.Id]; ok {
			productIdsMap[item.Product.Id] = append(productIdsMap[item.Product.Id], &shopCartItemsDetail{
				Item:       item,
				Promotions: []*models.PromotionItem{},
			})
		} else {
			productIdsMap[item.Product.Id] = []*shopCartItemsDetail{
				{Item: item, Promotions: []*models.PromotionItem{}},
			}
		}
	}
	details = map[string]*shopCartItemsDetail{}
	// 加载促销信息
	for productId, items := range productIdsMap {
		promotionItems := this.promotionRep.FindActivePromotionByProductId(ctx, productId)

		for _, item := range items {
			// 判断item是否属于活动
			item.Promotions = models.PromotionItems(promotionItems).IncludeItem(item.Item.GetID(), productId)

			if len(promotionItems) > 0 {
				if promotionItems[0].Promotion.Type == models.UnitSale {
					price := promotionItems[0].FindPriceByItemId(item.Item.GetID())
					if price != -1 {
						item.Item.PromotionPrice = price
					}
				}
			}
			details[item.Item.GetID()] = item
		}
	}
	return
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

func NewShopCartService(rep *repositories.ShopCartRep,
	itemRep *repositories.ItemRep,
	promotionRep *repositories.PromotionRep) *ShopCartService {
	return &ShopCartService{
		rep:          rep,
		itemRep:      itemRep,
		promotionRep: promotionRep,
	}
}
