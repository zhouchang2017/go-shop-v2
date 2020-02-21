package repositories

import (
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/models"
	"testing"
)

func TestProductRep_Create(t *testing.T) {
	newItems := make([]*models.Item,5)
	newItems[2] = &models.Item{
		Code:         "123",
		Product:      nil,
		Price:        0,
		OptionValues: nil,
		SalesQty:     0,
		Qty:          0,
		Inventories:  nil,
	}

	spew.Dump(newItems)
}
