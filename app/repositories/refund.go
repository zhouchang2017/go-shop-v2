package repositories

import (
	"context"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

type RefundRep struct {
	repository.IRepository
}

func (this *RefundRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "refund_order_no", Value: bsonx.Int64(-1)}}, // refund_order_no 唯一键
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
	}
}

func NewRefundRep(rep repository.IRepository) *RefundRep {
	repository := &RefundRep{rep}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Printf("model [%s] create indexs error:%s\n", repository.TableName(), err)
		panic(err)
	}
	return repository
}
