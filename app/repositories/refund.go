package repositories

import (
	"context"
	"go-shop-v2/app/models"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
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

// 通过退款单号查询订单
func (this *RefundRep) FindByRefundOrderNo(ctx context.Context, refundOrderNo string) (order *models.Refund, err error) {
	result := this.Collection().FindOne(ctx, bson.M{"refund_order_no": refundOrderNo})
	if result.Err() != nil {
		return nil, err2.Err404.F("refund order[%s] not fond", refundOrderNo)
	}
	order = &models.Refund{}
	if err := result.Decode(order); err != nil {
		return nil, err
	}
	return
}