package services

import (
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
)

func MakeBrandService() *BrandService {
	con := mongodb.GetConFn()
	rep := repositories.NewBrandRep(con)
	return NewBrandService(rep)
}

func MakeShopService() *ShopService {
	con := mongodb.GetConFn()
	return NewShopService(repositories.NewShopRep(con))
}

func MakeAdminService() *AdminService {
	con := mongodb.GetConFn()
	rep := repositories.NewAdminRep(con)
	return NewAdminService(rep)
}

func MakeItemService() *ItemService {
	rep := repositories.NewItemRep(mongodb.GetConFn())
	return NewItemService(rep)
}

func MakeProductService() *ProductService {
	rep := repositories.NewProductRep(mongodb.GetConFn())
	return NewProductService(rep)
}

func MakeInventoryService() *InventoryService {
	con := mongodb.GetConFn()
	rep := repositories.NewInventoryRep(con)
	historyRep := repositories.NewInventoryLogRep(con)
	return NewInventoryService(rep, historyRep, MakeShopService(), MakeProductService())
}

func MakeManualInventoryActionService() *ManualInventoryActionService {
	rep := repositories.NewManualInventoryActionRep(mongodb.GetConFn())
	return NewManualInventoryActionService(rep, MakeInventoryService(), MakeShopService(), MakeProductService())
}

func MakeCategoryService() *CategoryService {
	return NewCategoryService(repositories.NewCategoryRep(mongodb.GetConFn()))
}

func MakeArticleService() *ArticleService {
	return NewArticleService(repositories.NewArticleRep(mongodb.GetConFn()))
}

func MakeTopicService() *TopicService {
	return NewTopicService(repositories.NewTopicRep(mongodb.GetConFn()))
}

func MakeShopCartService() *ShopCartService {
	return NewShopCartService(repositories.NewShopCartRep(mongodb.GetConFn()))
}

func MakeUserService() *UserService {
	return NewUserService(repositories.NewUserRep(mongodb.GetConFn()))
}
