package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MongoModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func (this *MongoModel) SetID(id string) {
	ids, _ := primitive.ObjectIDFromHex(id)
	this.ID = ids
}
func (this *MongoModel) GetID() string {
	return this.ID.Hex()
}

func (this *MongoModel) GetCreatedAt() time.Time {
	return this.CreatedAt
}

func (this *MongoModel) GetUpdatedAt() time.Time {
	return this.UpdatedAt
}

func (this *MongoModel) SetCreatedAt(t time.Time) {
	this.CreatedAt = t
}

func (this *MongoModel) SetUpdatedAt(t time.Time) {
	this.UpdatedAt = t
}

