package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

func init() {
	register(NewOrderRep)
}

type OrderRep struct {
	*mongoRep
}

func NewOrderRep(con *mongodb.Connection) *OrderRep {
	return &OrderRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Order{}, con),
	}
}

