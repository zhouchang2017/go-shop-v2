package events

import "go-shop-v2/app/models"

type AdminUpdated struct {
	Admin *models.Admin
}
