package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"time"
)

// 促销
type PromotionService struct {
	rep *repositories.PromotionRep
}

func NewPromotionService(rep *repositories.PromotionRep) *PromotionService {
	return &PromotionService{rep: rep}
}

func (this *PromotionService) FindByIdWithItems(ctx context.Context, id string) (promotion *models.Promotion, err error) {
	return this.rep.FindByIdWithItems(ctx, id)
}

// 列表
func (this *PromotionService) Pagination(ctx context.Context, req *request.IndexRequest) (promotions []*models.Promotion, pagination response.Pagination, err error) {
	filters := req.Filters.Unmarshal()
	for key, value := range filters {
		req.AppendFilter(key, value)
	}
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	promotions = results.Result.([]*models.Promotion)
	pagination = results.Pagination
	return
}

func (this *PromotionService) FindById(ctx context.Context, id string) (promotion *models.Promotion, err error) {
	results := <-this.rep.FindById(ctx, id)
	if results.Error != nil {
		return nil, results.Error
	}
	return results.Result.(*models.Promotion), nil
}

type PromotionCreateOption struct {
	Name        string                  `json:"name"`                     // 促销活动名称
	Description string                  `json:"description"`              // 活动说明
	Items       []*models.PromotionItem `json:"items"`                    // 活动商品
	Type        int8                    `json:"type"`                     // 活动类型
	Mutex       bool                    `json:"mutex"`                    // 是否可叠加
	Rule        *models.PromotionRule   `json:"rule"`                     // 促销规则
	Policy      *models.PromotionPolicy `json:"policy"`                   // 促销策略
	Enable      bool                    `json:"enable"`                   // 是否开启
	BeginAt     time.Time               `json:"begin_at" form:"begin_at"` // 开始时间
	EndedAt     time.Time               `json:"ended_at" form:"ended_at"` // 结束时间
}

// 创建促销计划
func (this *PromotionService) Create(ctx context.Context, opt *PromotionCreateOption) (promotion *models.Promotion, err error) {
	model := &models.Promotion{
		Name:        opt.Name,
		Description: opt.Description,
		Items:       opt.Items,
		Type:        opt.Type,
		Mutex:       opt.Mutex,
		Rule:        opt.Rule,
		Policy:      opt.Policy,
		Enable:      opt.Enable,
		BeginAt:     opt.BeginAt,
		EndedAt:     opt.EndedAt,
	}
	created := <-this.rep.Create(ctx, model)
	if created.Error != nil {
		return nil, created.Error
	}
	return created.Result.(*models.Promotion), nil
}

// 更新促销计划
func (this *PromotionService) Update(ctx context.Context, model *models.Promotion, opt *PromotionCreateOption) (promotion *models.Promotion, err error) {
	model.Name = opt.Name
	model.Description = opt.Description
	model.Items = opt.Items
	model.Type = opt.Type
	model.Mutex = opt.Mutex
	model.Rule = opt.Rule
	model.Policy = opt.Policy
	model.Enable = opt.Enable
	model.BeginAt = opt.BeginAt
	model.EndedAt = opt.EndedAt
	saved := <-this.rep.Save(ctx, model)

	if saved.Error != nil {
		return nil, saved.Error
	}
	promotion = saved.Result.(*models.Promotion)

	return promotion, nil
}

// 删除
func (this *PromotionService) Delete(ctx context.Context, id string) (err error) {
	return <-this.rep.Delete(ctx, id)
}

// 获取Product显示价格
// 显示sku最低价格
// 如果产品,不存在促进价格 返回-1
func (this *PromotionService) FindProductsPrice(ctx context.Context, productIds ...string) (prices map[string]int64) {
	return this.rep.FindProductsPrice(ctx, productIds...)
}

// 获取Product对应的活动
func (this *PromotionService) FindActivePromotionByProductId(ctx context.Context, id string) (items []*models.PromotionItem) {
	return this.rep.FindActivePromotionByProductId(ctx, id)
}

// 计算订单促销价格
func (this *PromotionService) CalculateByOrder(ctx context.Context, items []*OrderItemCreateOption) *models.PromotionResult {
	// 组装数据
	var itemPromotions []*models.PromotionOrderItem
	for _, item := range items {
		itemPromotion := &models.PromotionOrderItem{
			ItemId:              item.ItemId,
			ProductId:           item.ProductId,
			Qty:                 item.Qty,
			UnMutexPromotions:   nil,
			MutexPromotion:      nil,
			MutexPromotionId:    item.MutexPromotion,
			UnMutexPromotionIds: item.UnMutexPromotions,
			Price:               item.Price,
		}

		var promotions []*models.PromotionItem
		// 互斥活动
		if item.MutexPromotion != nil {
			if promotion, err := this.rep.FindActivePromotion(ctx, *item.MutexPromotion, item.ProductId); err == nil {
				promotions = append(promotions, promotion)
			}
		}
		// 非互斥活动
		for _, id := range item.UnMutexPromotions {
			if promotion, err := this.rep.FindActivePromotion(ctx, id, item.ProductId); err == nil {
				promotions = append(promotions, promotion)
			}
		}
		var mutexPromotion *models.PromotionItem
		for _, promotion := range promotions {
			if promotion.Promotion.Mutex {
				if mutexPromotion == nil {
					itemPromotion.MutexPromotion = promotion
					mutexPromotion = promotion
				} else {
					// 出现两个互斥促销，异常
				}

			} else {
				itemPromotion.UnMutexPromotions = append(itemPromotion.UnMutexPromotions, promotion)
			}
		}

		itemPromotions = append(itemPromotions, itemPromotion)
	}
	return models.PromotionCalculate(itemPromotions...)
}
