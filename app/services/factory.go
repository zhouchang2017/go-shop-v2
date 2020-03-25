package services

import (
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
)

func MakeBrandService() *BrandService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Brand{}, mongodb.GetConFn())
	//brandCacheRep := repositories.NewRedisCache(&models.Brand{}, redis.GetConFn(), mongoRep)
	rep := repositories.NewBrandRep(mongoRep)
	return NewBrandService(rep)
}

func MakeShopService() *ShopService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Shop{}, mongodb.GetConFn())
	return NewShopService(repositories.NewShopRep(mongoRep))
}

func MakeAdminService() *AdminService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Admin{}, mongodb.GetConFn())
	rep := repositories.NewAdminRep(mongoRep)
	return NewAdminService(rep)
}

func newItemRep() *repositories.ItemRep {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Item{}, mongodb.GetConFn())
	//itemCacheRep := repositories.NewRedisCache(&models.Item{}, redis.GetConFn(), mongoRep)
	rep := repositories.NewItemRep(mongoRep)
	return rep
}

func newProductRep() *repositories.ProductRep {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Product{}, mongodb.GetConFn())
	//productCacheRep := repositories.NewRedisCache(&models.Product{}, redis.GetConFn(), mongoRep)
	return repositories.NewProductRep(mongoRep, newItemRep())
}

func MakeItemService() *ItemService {
	return NewItemService(newItemRep())
}

func MakeProductService() *ProductService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Product{}, mongodb.GetConFn())
	//productCacheRep := repositories.NewRedisCache(&models.Product{}, redis.GetConFn(), mongoRep)
	itemRep := newItemRep()
	rep := repositories.NewProductRep(mongoRep, itemRep)
	return NewProductService(rep, repositories.MakePromotionRep(), itemRep)
}

func MakeInventoryService() *InventoryService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Inventory{}, mongodb.GetConFn())
	rep := repositories.NewInventoryRep(mongoRep)

	historyRep := repositories.NewInventoryLogRep(repositories.NewBasicMongoRepositoryByDefault(&models.InventoryLog{}, mongodb.GetConFn()))

	return NewInventoryService(rep, historyRep, MakeShopService(), MakeProductService())
}

func MakeManualInventoryActionService() *ManualInventoryActionService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.ManualInventoryAction{}, mongodb.GetConFn())
	rep := repositories.NewManualInventoryActionRep(mongoRep)
	return NewManualInventoryActionService(rep, MakeInventoryService(), MakeShopService(), MakeProductService())
}

func MakeCategoryService() *CategoryService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Category{}, mongodb.GetConFn())
	rep := repositories.NewCategoryRep(mongoRep)
	return NewCategoryService(rep)
}

func MakeArticleService() *ArticleService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Article{}, mongodb.GetConFn())
	rep := repositories.NewArticleRep(mongoRep)
	return NewArticleService(rep)
}

func MakeTopicService() *TopicService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Topic{}, mongodb.GetConFn())
	rep := repositories.NewTopicRep(mongoRep)
	return NewTopicService(rep)
}

func MakeShopCartService() *ShopCartService {
	return NewShopCartService(
		repositories.MakeShopCartRep(),
		repositories.MakeItemRep(),
		repositories.MakePromotionRep(),
	)
}

func MakeUserService() *UserService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.User{}, mongodb.GetConFn())
	rep := repositories.NewUserRep(mongoRep)
	return NewUserService(rep)
}

func MakeBookmarkService() *BookmarkService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Bookmark{}, mongodb.GetConFn())
	rep := repositories.NewBookmarkRep(mongoRep)
	return NewBookmarkService(rep)
}

func MakeOrderService() *OrderService {
	return NewOrderService(
		repositories.MakeOrderRep(),
		MakePromotionService(),
		MakeProductService(),
	)
}

func MakeAddressService() *AddressService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.UserAddress{}, mongodb.GetConFn())
	rep := repositories.NewAddressRep(mongoRep)
	return NewAddressService(rep)
}

func MakePromotionItemService() *PromotionItemService {
	rep := repositories.NewPromotionItemRep(repositories.NewBasicMongoRepositoryByDefault(&models.PromotionItem{}, mongodb.GetConFn()))
	return NewPromotionItemService(rep, newProductRep())
}

func MakePromotionService() *PromotionService {
	return NewPromotionService(repositories.MakePromotionRep())
}

func MakePaymentService() *PaymentService {
	return NewPaymentService(repositories.MakePaymentRep(), repositories.MakeOrderRep())
}
