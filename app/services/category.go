package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type CategoryService struct {
	rep *repositories.CategoryRep
}

func NewCategoryService(rep *repositories.CategoryRep) *CategoryService {
	return &CategoryService{rep: rep}
}

// 列表
func (this *CategoryService) Pagination(ctx context.Context, req *request.IndexRequest) (categories []*models.Category, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Category), results.Pagination, nil
}

type CategoryCreateOption struct {
	Name    string               `json:"name" form:"name" binding:"required"`
}

type CategoryOptionForm struct {
	Id        string             `json:"id"`
	Name      string             `json:"name" form:"name"`
	Sort      int64              `json:"sort" form:"sort"`
	Thumbnail bool               `json:"thumbnail form:"thumbnail"`
	Values    []*optionValueForm `json:"values" form:"values"`
}

type optionValueForm struct {
	Code  string `json:"code" form:"code"`
	Value string `json:"value" form:"value"`
}

// 创建
func (this *CategoryService) Create(ctx context.Context, opt CategoryCreateOption) (category *models.Category, err error) {
	newCategory := models.NewCategory(opt.Name)

	//productOptions := []*models.ProductOption{}
	//for _, option := range opt.Options {
	//	productOption := models.NewProductOption(option.Name)
	//	productOption.Sort = option.Sort
	//	productOption.Thumbnail = option.Thumbnail
	//	values := []*models.OptionValue{}
	//	for _, value := range option.Values {
	//		values = append(values, productOption.NewValue(value.Value, value.Code))
	//	}
	//	productOption.Values = values
	//
	//	productOptions = append(productOptions, productOption)
	//}
	//
	//newCategory.Options = productOptions

	created := <-this.rep.Create(ctx, newCategory)
	if created.Error != nil {
		err = created.Error
		return
	}
	return created.Result.(*models.Category), nil
}

// 更新
func (this *CategoryService) Update(ctx context.Context, model *models.Category, opt CategoryCreateOption) (category *models.Category, err error) {
	model.Name = opt.Name

	//productOptions := []*models.ProductOption{}
	//for _, option := range opt.Options {
	//	var productOption *models.ProductOption
	//	if option.Id != "" {
	//		productOption = models.MakeProductOption(option.Id, option.Name, option.Sort)
	//	} else {
	//		productOption = models.NewProductOption(option.Name)
	//		productOption.Sort = option.Sort
	//	}
	//	productOption.Thumbnail = option.Thumbnail
	//
	//	values := []*models.OptionValue{}
	//	for _, value := range option.Values {
	//		values = append(values, productOption.NewValue(value.Value, value.Code))
	//	}
	//	productOption.Values = values
	//
	//	productOptions = append(productOptions, productOption)
	//}
	//model.Options = productOptions

	saved := <-this.rep.Save(ctx, model)

	if saved.Error != nil {
		err = saved.Error
		return
	}
	return saved.Result.(*models.Category), nil
}

// 详情
func (this *CategoryService) FindById(ctx context.Context, id string) (category *models.Category, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}
	return byId.Result.(*models.Category), nil
}

// 删除
func (this *CategoryService) Delete(ctx context.Context, id string) (err error) {
	return <-this.rep.Delete(ctx, id)
}

// 还原
func (this *CategoryService) Restore(ctx context.Context, id string) (category *models.Category, err error) {
	restored := <-this.rep.Restore(ctx, id)
	if restored.Error != nil {
		return nil, restored.Error
	}
	return restored.Result.(*models.Category), nil
}
