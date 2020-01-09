package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type ItemService struct {
	rep *repositories.ItemRep
}

func NewItemService(rep *repositories.ItemRep) *ItemService {
	return &ItemService{rep: rep}
}

// 列表
func (this *ItemService) Pagination(ctx context.Context, req *request.IndexRequest) (items []*models.Item, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Item), results.Pagination, nil
}

// 创建
func (this *ItemService) Create(ctx context.Context, model *models.Item) (item *models.Item, err error) {
	created := <-this.rep.Create(ctx, model)
	if created.Error != nil {
		err = created.Error
		return
	}
	return created.Result.(*models.Item), nil
}

// 更新
func (this *ItemService) Save(ctx context.Context, model *models.Item) (item *models.Item, err error) {
	saved := <-this.rep.Save(ctx, model)
	if saved.Error != nil {
		err = saved.Error
		return
	}
	return saved.Result.(*models.Item), nil
}

// 详情
func (this *ItemService) FindById(ctx context.Context, id string) (item *models.Item, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}
	return byId.Result.(*models.Item), nil
}

// 删除
func (this *ItemService) Delete(ctx context.Context, id string) (err error) {
	return <-this.rep.Delete(ctx, id)
}

// 还原
func (this *ItemService) Restore(ctx context.Context, id string) (item *models.Item, err error) {
	restored := <-this.rep.Restore(ctx, id)
	if restored.Error != nil {
		return nil, restored.Error
	}
	return restored.Result.(*models.Item), nil
}

func (this *ItemService) FindByIds(ctx context.Context, ids ...string) (items []*models.Item, err error) {
	results := <-this.rep.FindByIds(ctx, ids...)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Item), nil
}

func (this *ItemService) FindAssociatedByIds(ctx context.Context, ids ...string) (items []*models.AssociatedItem, err error) {
	results, err := this.FindByIds(ctx, ids...)
	if err != nil {
		return
	}

	for _, item := range results {
		items = append(items, item.ToAssociated())
	}
	return
}