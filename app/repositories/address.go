package repositories

import (
	"context"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressRep struct {
	repository.IRepository
}

func (this *AddressRep) SetIsDefault(ctx context.Context, filter interface{}, isDefault bool) (err error) {
	isDefaultValue := 0
	if isDefault {
		isDefaultValue = 1
	}
	_, err = this.Collection().UpdateMany(ctx, filter, bson.M{
		"$set": bson.M{
			"is_default": isDefaultValue,
		},
		"$currentDate": bson.M{
			"updated_at": true,
		},
	})
	return
}

func (this *AddressRep) SetIsDefaultByUserId(ctx context.Context, userId string, isDefault bool) (err error) {
	return this.SetIsDefault(ctx, bson.M{"user_id": userId}, isDefault)
}

func (this *AddressRep) SetIsDefaultById(ctx context.Context, id string, isDefault bool) (err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	return this.SetIsDefault(ctx, bson.M{"_id": objectID}, isDefault)
}

func NewAddressRep(IRepository repository.IRepository) *AddressRep {
	return &AddressRep{IRepository: IRepository}
}
