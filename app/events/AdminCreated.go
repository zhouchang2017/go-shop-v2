package events

import "go-shop-v2/app/models"

type AdminCreated struct {
	Admin *models.Admin
}
