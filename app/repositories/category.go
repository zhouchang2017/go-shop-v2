package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func init() {
	register(NewCategoryRep)
}

type CategoryRep struct {
	*mongoRep
}

// 添加销售属性
func (this *CategoryRep) AddOption(ctx context.Context, id string, opt *models.ProductOption) (err error) {
	ids, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", ids}}
	update := bson.M{
		"$addToSet": bson.D{{"options", opt}},
		"$currentDate": bson.D{
			{"updated_at", true},
		},
	}
	result, err := this.Collection().UpdateOne(ctx, filter, update)
	if result.ModifiedCount < 1 {
		log.Printf("[%s] add option modified count[%d] < 1 ,option[%+v]", id, result.ModifiedCount, opt)
	}
	return err
}

// 更新销售属性
func (this *CategoryRep) UpdateOption(ctx context.Context, id string, opt *models.ProductOption) (err error) {
	ids, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id":        ids,
		"options.id": opt.Id,
	}
	update := bson.M{
		"$set": bson.D{{"options.$", opt}},
		"$currentDate": bson.D{
			{"updated_at", true},
		},
	}
	result, err := this.Collection().UpdateOne(ctx, filter, update)
	if result.ModifiedCount < 1 {
		log.Printf("[%s] update option modified count[%d] < 1,option[%+v]", id, result.ModifiedCount, opt)
	}
	return err
}

// 删除销售属性
func (this *CategoryRep) DeleteOption(ctx context.Context, id string, optId string) (err error) {
	if id == "" {
		return err2.Err422.F("option id required!!")
	}
	if optId == "" {
		return err2.Err422.F("opt_id required!!")
	}
	ids, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err2.Err404
	}
	filter := bson.M{
		"_id": ids,
	}
	update := bson.M{
		"$pull": bson.D{{"options", bson.D{{"id", optId}}}},
	}
	result, err := this.Collection().UpdateOne(ctx, filter, update)
	if result.ModifiedCount < 1 {
		log.Printf("[%s] delete option modified count[%d] < 1,optionId[%s]", id, result.ModifiedCount, optId)
	}
	return err
}

func NewCategoryRep(con *mongodb.Connection) *CategoryRep {
	return &CategoryRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Category{}, con),
	}
}
