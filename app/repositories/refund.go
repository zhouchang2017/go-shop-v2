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

type RefundRep struct {
	repository.IRepository
}

func (this *RefundRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		// todo order_no order_id 加索引
		{
			Keys:    bsonx.Doc{{Key: "refund_no", Value: bsonx.Int64(-1)}}, // refund_no 唯一键
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
	}
}

func NewRefundRep(rep repository.IRepository) *RefundRep {
	repository := &RefundRep{rep}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Panicf("model [%s] create indexs error:%s\n", repository.TableName(), err)
	}
	return repository
}

func (this *RefundRep) FindMany(ctx context.Context, filter bson.M, opts ...*options.FindOptions) (refunds []*models.Refund, err error) {
	refunds = make([]*models.Refund, 0)
	cursor, err := this.Collection().Find(ctx, filter, opts...)
	if err != nil {
		return
	}

	if err := cursor.All(ctx, &refunds); err != nil {
		return refunds, err
	}
	return refunds, nil
}

func (this *RefundRep) FindByOrderNo(ctx context.Context, orderNo string) (refunds []*models.Refund, err error) {
	refunds = make([]*models.Refund, 0)
	cursor, err := this.Collection().Find(ctx, bson.M{"order_no": orderNo})
	if err != nil {
		return
	}

	if err := cursor.All(ctx, &refunds); err != nil {
		return refunds, err
	}
	return refunds, nil
}

// 通过退款单号查询订单
func (this *RefundRep) FindByRefundNo(ctx context.Context, refundNo string) (refund *models.Refund, err error) {
	result := this.Collection().FindOne(ctx, bson.M{"refund_no": refundNo})
	if result.Err() != nil {
		return nil, err2.Err404.F("refund_no[%s] not fond", refundNo)
	}
	refund = &models.Refund{}
	if err := result.Decode(refund); err != nil {
		return nil, err
	}
	return
}
