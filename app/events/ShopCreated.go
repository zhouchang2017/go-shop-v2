package events

import "go-shop-v2/app/models"

type ShopCreated struct {
	Shop *models.Shop
}
