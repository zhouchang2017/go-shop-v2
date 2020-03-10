package usecases

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"golang.org/x/sync/errgroup"
)

func PromotionWithItemsAndProduct(ctx context.Context, id string, promotionSrv *services.PromotionService, productService *services.ProductService) (promotion *models.Promotion, err error) {
	promotion, err = promotionSrv.FindByIdWithItems(ctx, id)
	if err != nil {
		return
	}

	var g errgroup.Group
	sem := make(chan struct{}, 10)
	items := make([]*models.PromotionItem, len(promotion.Items))
	for index, item := range promotion.Items {
		index, item := index, item
		sem <- struct{}{}
		g.Go(func() error {
			product, err := productService.FindById(ctx, item.ProductId)
			item.Product = product.ToAssociated()
			items[index] = item
			<-sem
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return promotion, err
	}
	promotion.Items = items
	return
}
