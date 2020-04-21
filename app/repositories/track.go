package repositories

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type TrackRep struct {
	repository.IRepository
}

func (this *TrackRep) FindByOrderNo(ctx context.Context, orderOn string) (traces []*models.Track, err error) {
	traces = make([]*models.Track, 0)
	cursor, err := this.Collection().Find(ctx, bson.M{
		"order_no": orderOn,
	})
	if err != nil {
		return traces, err
	}
	if err := cursor.All(ctx, &traces); err != nil {
		return traces, err
	}
	return traces, nil
}

func (this *TrackRep) FindOne(ctx context.Context, filter bson.M, opts ...*options.FindOneOptions) (trace *models.Track, err error) {
	one := this.Collection().FindOne(ctx, filter, opts...)
	if one.Err() != nil {
		return nil, err2.Err404
	}
	trace = &models.Track{}
	if err := one.Decode(trace); err != nil {
		return nil, err
	}
	return trace, nil
}

func (this *TrackRep) FindByWayBillId(ctx context.Context, deliveryId string, wayBillId string) (trace *models.Track, err error) {
	one := this.Collection().FindOne(ctx, bson.M{"delivery_id": deliveryId, "way_bill_id": wayBillId})
	if one.Err() != nil {
		return nil, err2.Err404
	}
	trace = &models.Track{}
	if err := one.Decode(trace); err != nil {
		return nil, err
	}
	return trace, nil
}

func (this *TrackRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "delivery_id", Value: bsonx.Int64(-1)},
				{Key: "way_bill_id", Value: bsonx.Int64(-1)},
			}, // order_no 唯一键
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
		{
			Keys: bsonx.Doc{
				{Key: "order_no", Value: bsonx.Int64(-1)},
			}, // order_no 唯一键
			Options: options.Index().SetBackground(true),
		},
	}
}

func NewTrackRep(IRepository repository.IRepository) *TrackRep {
	repository := &TrackRep{IRepository: IRepository}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Panicf("model [%s] create indexs error:%s\n", repository.TableName(), err)
	}
	return repository
}
