package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"golang.org/x/sync/errgroup"
	"sync"
)

type ItemService struct {
	rep *repositories.ItemRep
}

func NewItemService(rep *repositories.ItemRep) *ItemService {
	return &ItemService{rep: rep}
}

// 减库存
func (this *ItemService) DecQty(ctx context.Context, itemId string, qty int64) error {
	return this.rep.DecQty(ctx, itemId, qty)
}

// 加库存
func (this *ItemService) IncQty(ctx context.Context, itemId string, qty uint64) error {
	return this.rep.IncQty(ctx, itemId, qty)
}

func (this *ItemService) FindByProductId(ctx context.Context, productId string) (items []*models.Item) {
	items = []*models.Item{}
	results := <-this.rep.FindByProductId(ctx, productId)
	if results.Error != nil {
		return items
	}
	items = results.Result.([]*models.Item)
	return
}

// 列表
func (this *ItemService) Pagination(ctx context.Context, req *request.IndexRequest) (items []*models.Item, pagination response.Pagination, err error) {
	filters := req.Filters.Unmarshal()
	for key, value := range filters {
		req.AppendFilter(key, value)
	}
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	items = results.Result.([]*models.Item)

	return items, results.Pagination, nil
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

func (this *ItemService) Items(ctx context.Context, ids ...string) (items []*models.Item, err error) {
	var g errgroup.Group
	var mu sync.Mutex
	items = make([]*models.Item, len(ids))
	sem := make(chan struct{}, 10)
	for index, id := range ids {
		index, id := index, id
		sem <- struct{}{}
		g.Go(func() error {
			item, err := this.FindById(ctx, id)
			mu.Lock()
			items[index] = item
			mu.Unlock()
			<-sem
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return items, nil
}

// 详情
func (this *ItemService) FindById(ctx context.Context, id string) (item *models.Item, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}
	item = byId.Result.(*models.Item)
	item.Avatar = item.GetAvatar()
	return item, nil
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
	results := <-this.rep.FindByIds(ctx, ids)
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
