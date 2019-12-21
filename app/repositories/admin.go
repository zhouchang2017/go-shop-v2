package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

func init() {
	register(NewAdminRep)
}

type AdminRep struct {
	*mongoRep
}

func NewAdminRep(con *mongodb.Connection) *AdminRep {
	return &AdminRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Admin{}, con),
	}
}
