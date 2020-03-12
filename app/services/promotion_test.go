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

	items := service.FindActivePromotionByProductId(context.Background(), "5e579319b47683085db96caf")
	spew.Dump(items)
}
