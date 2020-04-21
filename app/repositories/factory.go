package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

func MakeAdminRep() *AdminRep {
	rep := NewBasicMongoRepositoryByDefault(&models.Admin{}, mongodb.GetConFn())
	return NewAdminRep(rep)
}

func MakeProductRep() *ProductRep {
	rep := NewBasicMongoRepositoryByDefault(&models.Product{}, mongodb.GetConFn())
	return NewProductRep(rep, MakeItemRep())
}

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

func MakeRefundRep() *RefundRep {
	refundMongoRep := NewBasicMongoRepositoryByDefault(&models.Refund{}, mongodb.GetConFn())
	return NewRefundRep(refundMongoRep)
}


func MakeCommentRep() *CommentRep {
	rep := NewBasicMongoRepositoryByDefault(&models.Comment{}, mongodb.GetConFn())
	return NewCommentRep(rep)
}

func MakePaymentRep() *PaymentRep {
	rep := NewBasicMongoRepositoryByDefault(&models.Payment{}, mongodb.GetConFn())
	return NewPaymentRep(rep)
}

func MakeTrackRep() *TrackRep {
	rep := NewBasicMongoRepositoryByDefault(&models.Track{}, mongodb.GetConFn())
	return NewTrackRep(rep)
}
