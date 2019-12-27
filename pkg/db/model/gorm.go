package model

import (
	"strconv"
	"time"
)

type GORMModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

func (this *GORMModel) GetID() string {
	return strconv.Itoa(int(this.ID))
}

func (this *GORMModel) GetCreatedAt() time.Time {
	return this.CreatedAt
}

func (this *GORMModel) GetUpdatedAt() time.Time {
	return this.UpdatedAt
}

func (this *GORMModel) SetCreatedAt(t time.Time) {
	this.CreatedAt = t
}

func (this *GORMModel) SetUpdatedAt(t time.Time) {
	this.UpdatedAt = t
}

func (this *GORMModel) IsSoftDeleted() bool {
	return this.DeletedAt != nil
}
