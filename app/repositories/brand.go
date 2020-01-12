package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)


type BrandRep struct {
	*mongoRep
}

func NewBrandRep(con *mongodb.Connection) *BrandRep {
	return &BrandRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Brand{}, con),
	}
}
