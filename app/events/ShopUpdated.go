package events

import (
	"go-shop-v2/app/models"
)

type ShopUpdated struct {
	Shop *models.Shop
}