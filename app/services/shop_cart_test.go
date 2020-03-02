package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/utils"
	"testing"
)

func TestShopCartService_Add(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	itemService := MakeItemService()
	shopCartService := MakeShopCartService()
	req := &request.IndexRequest{}
	req.Page = 1

	items, _, err := itemService.Pagination(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range items {
		num := utils.RandomInt(5)
		if num > 2 {
			err := shopCartService.Add(context.Background(), "5e5351eeb624a18a80352c7c", item.GetID(), num)
			if err != nil {
				t.Fatal(err)
			}
		}

	}
}

func TestShopCartService_Index(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	shopCartService := MakeShopCartService()

	var page int64 = 1
	var perPage int64 = 15
	var hasNextPage = true
	for hasNextPage {
		_, pagination, err := shopCartService.Index(context.Background(), "5e5351eeb624a18a80352c7c", page, perPage)
		if err != nil {
			t.Fatal(err)
		}
		hasNextPage = pagination.HasNextPage
		page++

		spew.Dump(pagination)
	}

}

func TestShopCartService_Delete(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	shopCartService := MakeShopCartService()

	count := shopCartService.Count(context.Background(), "5e5351eeb624a18a80352c7c")

	var deleteCount int64 = 15

	items, _, err := shopCartService.Index(context.Background(), "5e5351eeb624a18a80352c7c", 1, deleteCount)

	var ids []string
	for _, item := range items {
		ids = append(ids, item.ItemId)
	}
	err = shopCartService.Delete(context.Background(), "5e5351eeb624a18a80352c7c", ids...)
	if err != nil {
		t.Fatal(err)
	}

	deleteNum := count - shopCartService.Count(context.Background(), "5e5351eeb624a18a80352c7c")

	if deleteNum != 15 {
		t.Errorf("预计删除数量与实际不符，预计删除%d,实际删除%d", deleteCount, deleteNum)
	}

}

func TestShopCartService_Toggle(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	shopCartService := MakeShopCartService()

	var toggleCount int64 = 15
	var updateChecked = true
	items, _, err := shopCartService.Index(context.Background(), "5e5351eeb624a18a80352c7c", 1, toggleCount)

	var ids []string
	for _, item := range items {
		ids = append(ids, item.ItemId)

	}
	err = shopCartService.Toggle(context.Background(), "5e5351eeb624a18a80352c7c", updateChecked, ids...)
	if err != nil {
		t.Fatal(err)
	}

	index, _, _ := shopCartService.Index(context.Background(), "5e5351eeb624a18a80352c7c", 1, toggleCount)

	for _, item := range index {
		if item.Checked != updateChecked {
			t.Errorf("预计选择状态与实际不符，预计选中状态%t,实际选中状态%t", updateChecked, item.Checked)
		}
	}

}

func TestShopCartService_UpdateQty(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	shopCartService := MakeShopCartService()

	err := shopCartService.UpdateQty(context.Background(), "5e5351eeb624a18a80352c7c", "5e57b4d2b47683085db96cef", 5)
	if err != nil {
		t.Fatal(err)
	}

}

func TestShopCartService_Update(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	shopCartService := MakeShopCartService()

	update, err := shopCartService.Update(context.Background(), "5e5351eeb624a18a80352c7c", "5e57b4d2b47683085db96cef", "5e57b4d2b47683085db96cee", 8)
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(update)
}
