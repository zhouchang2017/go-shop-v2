package services

import (
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/db/mongodb"
)

func MakeBrandService() *BrandService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Brand{}, mongodb.GetConFn())
	brandCacheRep := repositories.NewRedisCache(&models.Brand{}, redis.GetConFn(), mongoRep)
	rep := repositories.NewBrandRep(brandCacheRep)
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
	itemCacheRep := repositories.NewRedisCache(&models.Item{}, redis.GetConFn(), mongoRep)
	rep := repositories.NewItemRep(itemCacheRep)
	return rep
}

func MakeItemService() *ItemService {
	return NewItemService(newItemRep())
}

func MakeProductService() *ProductService {
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Product{}, mongodb.GetConFn())
	productCacheRep := repositories.NewRedisCache(&models.Product{}, redis.GetConFn(), mongoRep)
	rep := repositories.NewProductRep(productCacheRep, newItemRep())
	return NewProductService(rep)
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
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.ShopCart{}, mongodb.GetConFn())
	rep := repositories.NewShopCartRep(mongoRep)
	return NewShopCartService(rep)
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
	mongoRep := repositories.NewBasicMongoRepositoryByDefault(&models.Order{}, mongodb.GetConFn())
	rep := repositories.NewOrderRep(mongoRep)
	return NewOrderService(rep)
}
