package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"sync"
	"testing"
)

func TestBookmarkService_Add(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeBookmarkService()

	productService := MakeProductService()

	req := &request.IndexRequest{}
	req.Page = -1
	data, _, err := productService.Pagination(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, product := range data {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := service.Add(context.Background(), "5e5351eeb624a18a80352c7c", id)
			if err != nil {
				t.Fatal(err)
			}
		}(product.GetID())
	}
	wg.Wait()

}

func TestBookmarkService_Count(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeBookmarkService()

	count := service.Count(context.Background(), "5e5351eeb624a18a80352c7c")
	spew.Dump(count)
}

func TestBookmarkService_Index(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeBookmarkService()

	ids, pagination, err := service.Index(context.Background(), "5e5351eeb624a18a80352c7c", 2, 15)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(ids)
	spew.Dump(pagination)
}
