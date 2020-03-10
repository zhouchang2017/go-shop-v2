package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"testing"
)

func TestPromotionItemService_Pagination(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()


	service:= MakePromotionItemService()

	data, _, err := service.Pagination(context.Background(), &request.IndexRequest{})
	if err!=nil {
		t.Fatal(err)
	}
	spew.Dump(data)
}
