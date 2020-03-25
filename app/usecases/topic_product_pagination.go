package usecases

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/response"
)

func TopicProductPagination(ctx context.Context, topicId string, page int64, perPage int64, topicSrv *services.TopicService, productSrv *services.ProductService) (products []*models.Product, pagination response.Pagination, err error) {
	ids, pagination, err := topicSrv.ProductsPagination(ctx, topicId, page, perPage)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}
	products, err = productSrv.FindByIds(ctx, ids)
	if err != nil {
		return
	}
	return
}
