package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type BookmarkRep struct {
	repository.IRepository
}

// 分页
func (this *BookmarkRep) Index(ctx context.Context, userId string, page int64, perPage int64) (ids []string, pagination response.Pagination, err error) {
	ids = []string{}
	count := this.Count(ctx, userId)
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = 15
	}
	skip := (page - 1) * 15

	pagination = response.Pagination{
		Total:       count,
		CurrentPage: page,
		PerPage:     perPage,
		HasNextPage: page*perPage < count,
	}

	result := this.Collection().FindOne(ctx, bson.M{"user_id": userId}, options.FindOne().SetProjection(bson.M{
		"product_ids": bson.M{"$slice": bson.A{skip, perPage * page}},
	}))

	if result.Err() != nil {
		err = result.Err()
		return
	}
	var res models.Bookmark
	if err := result.Decode(&res); err != nil {
		return []string{}, pagination, err
	}
	for _, id := range res.ProductIds {
		ids = append(ids, id)
	}
	return ids, pagination, nil
}

// 添加
func (this *BookmarkRep) Add(ctx context.Context, userId string, productId string) (err error) {
	_, err = this.Collection().UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{
		"$addToSet": bson.M{
			"product_ids": productId,
		},
		"$currentDate": bson.M{
			"updated_at": true,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}, options.Update().SetUpsert(true))
	return
}

// 移除
func (this *BookmarkRep) Remove(ctx context.Context, userId string, productIds ...string) (err error) {
	_, err = this.Collection().UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{
		"$pull": bson.M{
			"product_ids": bson.M{"$in": productIds},
		},
		"$currentDate": bson.M{
			"updated_at": true,
		},
	})
	return err
}

type countRep struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Count int64              `json:"count"`
}

// 总数
func (this *BookmarkRep) Count(ctx context.Context, userId string) (count int64) {
	aggregate, err := this.Collection().Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match", bson.M{"user_id": userId}}},
		bson.D{{"$project", bson.M{"count": bson.M{"$size": "$product_ids"}}}},
	})
	if err != nil {
		return 0
	}
	var res []countRep
	err = aggregate.All(ctx, &res)
	if err != nil {
		return 0
	}
	if len(res) > 0 {
		return res[0].Count
	}
	return 0
}

func (this *BookmarkRep) FindByProductId(ctx context.Context, userId string, productId string) (bookmark *models.Bookmark) {
	result := this.Collection().FindOne(ctx, bson.M{"user_id": userId, "product_ids": bson.M{"$in": bson.A{productId}}}, options.FindOne().SetProjection(bson.M{"product_ids": 0}))
	if result.Err() != nil {
		return nil
	}
	bookmark = &models.Bookmark{}
	if err := result.Decode(bookmark); err != nil {
		return nil
	}
	return bookmark
}

func NewBookmarkRep(rep repository.IRepository) *BookmarkRep {
	return &BookmarkRep{rep}
}
