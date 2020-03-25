package repositories

import (
	"context"
	"go-shop-v2/app/models"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TopicRep struct {
	repository.IRepository
}

// 产品总数
func (this *TopicRep) ProductCount(ctx context.Context, id string) (count int64) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0
	}
	aggregate, err := this.Collection().Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match", bson.M{"_id": objId}}},
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

// 分页
func (this *TopicRep) Products(ctx context.Context, topicId string, page int64, perPage int64) (ids []string, pagination response.Pagination, err error) {
	ids = []string{}
	count := this.ProductCount(ctx, topicId)
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

	objId, err := primitive.ObjectIDFromHex(topicId)
	if err != nil {
		err = err2.Err404
		return
	}
	result := this.Collection().FindOne(ctx, bson.M{"_id": objId}, options.FindOne().SetProjection(bson.M{
		"product_ids": bson.M{"$slice": bson.A{skip, perPage * page}},
	}))

	if result.Err() != nil {
		err = result.Err()
		return
	}
	var res models.Topic
	if err := result.Decode(&res); err != nil {
		return []string{}, pagination, err
	}
	for _, id := range res.ProductIds {
		ids = append(ids, id)
	}
	return ids, pagination, nil
}

func NewTopicRep(rep repository.IRepository) *TopicRep {
	return &TopicRep{rep}
}
