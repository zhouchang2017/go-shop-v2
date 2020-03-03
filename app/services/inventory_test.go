package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestInventoryService_Put(t *testing.T) {
	//mongodb.TestConnect()
	//defer mongodb.Close()
	//
	//con := mongodb.GetConFn()
	//
	//shopRep := repositories.NewShopRep(con)
	//itemRep := repositories.NewItemRep(con)
	////rep := repositories.NewManualInventoryActionRep(con)
	//service := MakeInventoryService()
	//
	//type TestData struct {
	//	ShopId string
	//	ItemId string
	//	Qty    int64
	//	Status int8
	//}
	//
	//allItemRes := <-itemRep.FindAll(context.Background())
	//if allItemRes.Error != nil {
	//	t.Fatal(allItemRes.Error)
	//}
	//items := allItemRes.Result.([]*models.Item)
	//allShopRes := <-shopRep.FindAll(context.Background())
	//if allShopRes.Error != nil {
	//	t.Fatal(allShopRes.Error)
	//}
	//shops := allShopRes.Result.([]*models.Shop)
	//var data []TestData
	//for _, item := range items {
	//	for _, shop := range shops {
	//		data = append(data, TestData{
	//			ShopId: shop.GetID(),
	//			ItemId: item.GetID(),
	//			Qty:    utils.RandomInt(10),
	//			Status: int8(utils.RandomInt(3)),
	//		})
	//	}
	//}
	//
	//for _, d := range data {
	//	inventory, err := service.Put(context.Background(),
	//		d.ShopId,
	//		d.ItemId,
	//		d.Qty,
	//		d.Status)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	t.Logf("puted inventory %+v\n", inventory)
	//}

}

func TestInventoryService_Take(t *testing.T) {
	//mongodb.TestConnect()
	//defer mongodb.Close()
	//
	//service := MakeInventoryService()
	//inventory, err := service.Take(context.Background(), "5dfcd13c540b36a9ac259c65", 8)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Logf("taked inventory %+v\n", inventory)
}

func TestInventoryService_GetRepository(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeInventoryService()

	res := &request.IndexRequest{}
	es, pagination, err := service.Aggregate(context.Background(), res)
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(es)

	spew.Dump(pagination)
}

func TestInventoryService_Search(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeInventoryService()

	inventories, err := service.Search(context.Background(), &repositories.QueryOption{
		ProductId: "5e33fcd21f50bb25b7060ec0",
		Qty:       bson.M{"$gt": 5},
		Status:    []int8{0},
	})

	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(inventories)
}

func TestInventoryService_SearchItemsQty(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeInventoryService()

	qty := service.SearchItemsQty(context.Background(), "5e51e253ecbe820cbd5f6d7f","5e51e253ecbe820cbd5f6d74","5e577e370d3f4744961cfd0f")
	spew.Dump(qty)
}