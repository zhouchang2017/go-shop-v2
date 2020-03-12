package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"golang.org/x/sync/errgroup"
)

type PromotionItemService struct {
	rep        *repositories.PromotionItemRep
	productRep *repositories.ProductRep
}

func NewPromotionItemService(rep *repositories.PromotionItemRep, productRep *repositories.ProductRep) *PromotionItemService {
	return &PromotionItemService{rep: rep, productRep: productRep}
}

func (this *PromotionItemService) Pagination(ctx context.Context, req *request.IndexRequest) (promotions []*models.PromotionItem, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	promotions = results.Result.([]*models.PromotionItem)
	pagination = results.Pagination

	var g errgroup.Group
	res := []*models.PromotionItem{}
	sem := make(chan struct{}, 10)
	for _, promotion := range promotions {
		promotion := promotion
		sem <- struct{}{}
		g.Go(func() error {
			product, err := this.productRep.WithItems(ctx, promotion.ProductId)
			if err != nil {
				if err == err2.Err404 {
					<-sem
					return nil
				}
			}
			promotion.Product = product.ToAssociated()

			for _, unit := range promotion.Units {
				for _, item := range product.Items {
					if unit.ItemId == item.GetID() {
						unit.Item = item.ToAssociated()
					}
				}
			}
			res = append(res, promotion)
			<-sem
			return err
		})
	}
	if err := g.Wait(); err != nil {

		return res, pagination, err
	}
	promotions = res
	return
}
