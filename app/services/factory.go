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
	return NewShopService(repositories.NewShopRep(con), repositories.NewAdminRep(con))
}

func MakeAdminService() *AdminService {
	con := mongodb.GetConFn()
	rep := repositories.NewAdminRep(con)
	return NewAdminService(rep, MakeShopService())
}

func MakeProductService() *ProductService {
	rep := repositories.NewProductRep(mongodb.GetConFn())
	return NewProductService(rep)
}

func MakeInventoryService() *InventoryService {
	con := mongodb.GetConFn()
	rep := repositories.NewInventoryRep(con)
	return NewInventoryService(rep, MakeShopService(), MakeProductService())
}

func MakeManualInventoryActionService() *ManualInventoryActionService {
	rep := repositories.NewManualInventoryActionRep(mongodb.GetConFn())
	return NewManualInventoryActionService(rep, MakeInventoryService(), MakeShopService(), MakeProductService())
}

func MakeCategoryService() *CategoryService  {
	return NewCategoryService(repositories.NewCategoryRep(mongodb.GetConFn()))
}