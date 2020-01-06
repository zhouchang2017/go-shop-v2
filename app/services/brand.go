package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type BrandService struct {
	rep *repositories.BrandRep
}

func NewBrandService(rep *repositories.BrandRep) *BrandService {
	return &BrandService{rep: rep}
}

// 列表
func (this *BrandService) Pagination(ctx context.Context, req *request.IndexRequest) (brands []*models.Brand, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Brand), results.Pagination, nil
}

// 创建
func (this *BrandService) Create(ctx context.Context, name string) (brand *models.Brand, err error) {
	created := <-this.rep.Create(ctx, &models.Brand{Name: name})
	if created.Error != nil {
		err = created.Error
		return
	}
	return created.Result.(*models.Brand), nil
}

// 更新
func (this *BrandService) Update(ctx context.Context, model *models.Brand, name string) (brand *models.Brand, err error) {
	if model.Name != name {
		model.Name = name
		saved := <-this.rep.Save(ctx, model)
		if saved.Error != nil {
			err = saved.Error
			return
		}
		brand = saved.Result.(*models.Brand)
	}
	return model, nil
}

// 详情
func (this *BrandService) FindById(ctx context.Context, id string) (brand *models.Brand, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}
	return byId.Result.(*models.Brand), nil
}

// 删除
func (this *BrandService) Delete(ctx context.Context, id string) (err error) {
	return <-this.rep.Delete(ctx, id)
}

// 还原
func (this *BrandService) Restore(ctx context.Context, id string) (brand *models.Brand, err error) {
	restored := <-this.rep.Restore(ctx, id)
	if restored.Error != nil {
		return nil, restored.Error
	}
	return restored.Result.(*models.Brand), nil
}
