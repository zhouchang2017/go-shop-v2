package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

type FailedRefundRep struct {
	repository.IRepository
}

func (this *FailedRefundRep) Write(ctx context.Context, failedRefund *models.FailedRefund) error {
	_, err := this.Collection().UpdateOne(ctx, bson.M{"refund_on": failedRefund.RefundOn}, bson.M{
		"$set": bson.M{
			"order_no":     failedRefund.OrderNo,
			"err_code":     failedRefund.ErrCode,
			"err_code_des": failedRefund.ErrCodeDes,
		},
		"$currentDate": bson.M{
			"updated_at": true,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}, options.Update().SetUpsert(true))
	return err
}

func (this *FailedRefundRep) ClearFailedByRefundOn(ctx context.Context, id string) error {
	_, err := this.Collection().DeleteOne(ctx, bson.M{
		"refund_on": id,
	})
	return err
}

func (this *FailedRefundRep) FindByRefundOn(ctx context.Context, refundOn string) (failedRefund *models.FailedRefund, err error) {
	one := this.Collection().FindOne(ctx, bson.M{
		"refund_on": refundOn,
	})
	if one.Err() != nil {
		return nil, one.Err()
	}
	failedRefund = &models.FailedRefund{}
	if err := one.Decode(failedRefund); err != nil {
		return nil, err
	}
	return failedRefund, nil
}
func (this *FailedRefundRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "order_no", Value: bsonx.Int64(-1)}},
			Options: options.Index().SetBackground(true),
		},
		{
			Keys:    bsonx.Doc{{Key: "refund_on", Value: bsonx.Int64(-1)}},
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
	}
}

func NewFailedRefundRep(IRepository repository.IRepository) *FailedRefundRep {
	repository := &FailedRefundRep{IRepository: IRepository}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Printf("model [%s] create indexs error:%s\n", repository.TableName(), err)
		panic(err)
	}
	return repository
}
