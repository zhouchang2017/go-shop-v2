package model

import "time"

type IModel interface {
	GetID() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreatedAt(t time.Time)
	SetUpdatedAt(t time.Time)
	IsSoftDeleted() bool
}
