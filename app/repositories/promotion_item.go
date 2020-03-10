package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PromotionItemRep struct {
	repository.IRepository
}

func (this *PromotionItemRep) FindByPromotionId(ctx context.Context, promotionId string, opts ...*options.FindOptions) (items []*models.PromotionItem, err error) {
	items = []*models.PromotionItem{}
	find, err := this.Collection().Find(ctx, bson.M{"promotion.id": promotionId}, opts...)
	if err != nil {
		return nil, err
	}
	if err := find.All(ctx, &items); err != nil {
		return items, err
	}
	return items, nil
}

func NewPromotionItemRep(IRepository repository.IRepository) *PromotionItemRep {
	return &PromotionItemRep{IRepository: IRepository}
}
