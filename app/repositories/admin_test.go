package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	"testing"
)

func TestNewAdminRep(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	rep := NewAdminRep(mongodb.GetConFn())

	admin := &models.Admin{
		Username: "zhouchang",
		Nickname: "小周",
		Type:     "root",
		Shops:    []*models.AssociatedShop{},
	}
	admin.SetPassword("123456")
	created := <-rep.Create(context.Background(), admin)

	if created.Error != nil {
		t.Fatal(created.Error)
	}
	t.Logf("%+v", created.Result)
}
