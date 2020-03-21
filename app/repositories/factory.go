package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

func MakePromotionRep() *PromotionRep {
	promotionItemRep := NewPromotionItemRep(NewBasicMongoRepositoryByDefault(&models.PromotionItem{}, mongodb.GetConFn()))
	return NewPromotionRep(NewBasicMongoRepositoryByDefault(&models.Promotion{}, mongodb.GetConFn()), promotionItemRep)
}

func MakeItemRep() *ItemRep {
	rep := NewBasicMongoRepositoryByDefault(&models.Item{}, mongodb.GetConFn())
	return NewItemRep(rep)
}

func MakeShopCartRep() *ShopCartRep {
	rep := NewBasicMongoRepositoryByDefault(&models.ShopCart{}, mongodb.GetConFn())
	return NewShopCartRep(rep)
}

func MakeOrderRep() *OrderRep {
	orderMongoRep := NewBasicMongoRepositoryByDefault(&models.Order{}, mongodb.GetConFn())
	return NewOrderRep(orderMongoRep)
}

func MakeCommentRep() *CommentRep {
	rep := NewBasicMongoRepositoryByDefault(&models.Comment{}, mongodb.GetConFn())
	return NewCommentRep(rep)
}
