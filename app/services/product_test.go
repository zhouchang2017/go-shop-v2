package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"testing"
)

func TestProductService_List(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeProductService()

	products, pagination, err := service.List(context.Background(), &request.IndexRequest{})

	spew.Dump(products)
	spew.Dump(pagination)

	if err!=nil {
		t.Fatal(err)
	}

}
