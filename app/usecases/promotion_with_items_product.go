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
	items := make([]*models.PromotionItem, 0)
	for _, item := range promotion.Items {
		item := item
		sem <- struct{}{}
		g.Go(func() error {
			product, err := productService.FindById(ctx, item.ProductId)
			if err == nil {
				item.Product = product.ToAssociated()
				items = append(items, item)
			}
			<-sem
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return promotion, err
	}
	promotion.Items = items
	return
}
