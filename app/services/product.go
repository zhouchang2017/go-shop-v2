package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
)

func init() {
	register(NewProductService)
}

type ProductService struct {
	rep *repositories.ProductRep
}

func NewProductService(rep *repositories.ProductRep) *ProductService {
	return &ProductService{rep: rep}
}

func (p *ProductService) FindByIdWithItems(ctx context.Context, id string) (product *models.Product, err error) {
	res := <-p.rep.FindById(ctx, id)
	if res.Error != nil {
		return nil, res.Error
	}
	product = res.Result.(*models.Product)
	itemRes := <-p.rep.GetItemRep().FindByProductId(ctx, id)
	if itemRes.Error != nil {
		return nil, itemRes.Error
	}
	product.Items = itemRes.Result.([]*models.Item)
	return product, nil
}
