package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type ProductService struct {
	rep         *repositories.ProductRep
	ItemService *ItemService
}

func NewProductService(rep *repositories.ProductRep) *ProductService {
	return &ProductService{
		rep:         rep,
		ItemService: NewItemService(rep.GetItemRep()),
	}
}

func (this *ProductService) FindItemById(ctx context.Context, id string) (item *models.Item, err error) {
	return this.rep.FindItemById(ctx, id)
}

func (this *ProductService) FindByIdWithItems(ctx context.Context, id string) (product *models.Product, err error) {
	res := <-this.rep.FindById(ctx, id)
	if res.Error != nil {
		return nil, res.Error
	}
	product = res.Result.(*models.Product)
	itemRes := <-this.rep.GetItemRep().FindByProductId(ctx, id)
	if itemRes.Error != nil {
		return nil, itemRes.Error
	}
	product.Items = itemRes.Result.([]*models.Item)
	return product, nil
}

// 列表
func (this *ProductService) Pagination(ctx context.Context, req *request.IndexRequest) (products []*models.Product, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Product), results.Pagination, nil
}

type ProductCreateOption struct {
	Name         string                     `json:"name" form:"name" binding:"required,max=255"`
	Code         string                     `json:"code" form:"code" binding:"required,max=255"`
	Brand        *models.AssociatedBrand    `json:"brand" form:"brand"`
	Category     *models.AssociatedCategory `json:"category" form:"category"`
	Attributes   []*models.ProductAttribute `json:"attributes" form:"attributes"`
	Options      []*models.ProductOption    `json:"options" form:"options"`
	Items        []*models.Item             `json:"items"`
	Description  string                     `json:"description"`
	Price        int64                      `json:"price"`
	FakeSalesQty int64                      `json:"fake_sales_qty" form:"fake_sales_qty"`
	Images       []*qiniu.Resource          `json:"images" form:"images"`
	OnSale       bool                       `json:"on_sale" form:"on_sale"`
}

// 创建
func (this *ProductService) Create(ctx context.Context, opt ProductCreateOption) (product *models.Product, err error) {
	model := &models.Product{
		Name:         opt.Name,
		Code:         opt.Code,
		Brand:        opt.Brand,
		Category:     opt.Category,
		Options:      opt.Options,
		Attributes:   opt.Attributes,
		Description:  opt.Description,
		Price:        opt.Price,
		FakeSalesQty: opt.FakeSalesQty,
		Images:       opt.Images,
		OnSale:       opt.OnSale,
		Items:        opt.Items,
	}

	created := <-this.rep.Create(ctx, model)
	if created.Error != nil {
		err = created.Error
		return
	}
	return created.Result.(*models.Product), nil
}

// 更新
func (this *ProductService) Update(ctx context.Context, model *models.Product, opt ProductCreateOption) (product *models.Product, err error) {
	model.Name = opt.Name
	model.Brand = opt.Brand
	model.Items = opt.Items
	model.Options = opt.Options
	model.Attributes = opt.Attributes
	model.Description = opt.Description
	model.Price = opt.Price
	model.FakeSalesQty = opt.FakeSalesQty
	model.Images = opt.Images
	model.OnSale = opt.OnSale

	saved := <-this.rep.Save(ctx, model)

	if saved.Error != nil {
		return nil, saved.Error
	}
	return saved.Result.(*models.Product), nil
}

// 详情
func (this *ProductService) FindById(ctx context.Context, id string) (product *models.Product, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}
	return byId.Result.(*models.Product), nil
}

// 删除
func (this *ProductService) Delete(ctx context.Context, id string) (err error) {
	return <-this.rep.Delete(ctx, id)
}

// 还原
func (this *ProductService) Restore(ctx context.Context, id string) (product *models.Product, err error) {
	restored := <-this.rep.Restore(ctx, id)
	if restored.Error != nil {
		return nil, restored.Error
	}
	return restored.Result.(*models.Product), nil
}
