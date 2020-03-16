package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/db/mongodb"
	"testing"
)

func TestPromotionService_FindActivePromotionByProductId(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakePromotionService()

	items := service.FindActivePromotionByProductId(context.Background(), "5e69b789d9acdd33dafb742e")
	spew.Dump(items)
}

func TestPromotionService_FindActivePromotionUnitSaleByProductId(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakePromotionService()
	items := service.FindActivePromotionByProductId(context.Background(), "5e69b789d9acdd33dafb742e")
	spew.Dump(items)

}