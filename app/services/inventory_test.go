package services

import (
	"context"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"testing"
)

func TestInventoryService_Put(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	con := mongodb.GetConFn()
	service := NewInventoryService(
		repositories.NewInventoryRep(con),
		repositories.NewShopRep(con),
		repositories.NewItemRep(con))

	type TestData struct {
		ShopId string
		ItemId string
		Qty    int64
		Status int8
	}

	data := []TestData{
		{
			ShopId: "5dfba56657a5dd253b675221",
			ItemId: "5dfb3e5248c83255b822f5ff",
			Qty:    10,
			Status: 2,
		},
		{
			ShopId: "5dfba56657a5dd253b675221",
			ItemId: "5dfb3e5248c83255b822f603",
			Qty:    1,
			Status: 0,
		},
		{
			ShopId: "5dfba56657a5dd253b675221",
			ItemId: "5dfb3e5248c83255b822f603",
			Qty:    5,
			Status: 2,
		},
		{
			ShopId: "5dfb759ebfceae2e2b857eca",
			ItemId: "5dfb3e5248c83255b822f603",
			Qty:    2,
			Status: 2,
		},
		{
			ShopId: "5dfba56657a5dd253b675221",
			ItemId: "5dfb3e5248c83255b822f5ff",
			Qty:    1,
			Status: 2,
		},
	}

	for _, d := range data {
		inventory, err := service.Put(context.Background(),
			d.ShopId,
			d.ItemId,
			d.Qty,
			d.Status)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("puted inventory %+v\n", inventory)
	}

}

func TestInventoryService_Take(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	con := mongodb.GetConFn()
	service := NewInventoryService(
		repositories.NewInventoryRep(con),
		repositories.NewShopRep(con),
		repositories.NewItemRep(con))
	inventory, err := service.Take(context.Background(), "5dfcd13c540b36a9ac259c65", 8)
	if err!=nil {
		t.Fatal(err)
	}

	t.Logf("taked inventory %+v\n", inventory)
}
