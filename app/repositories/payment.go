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

type PaymentRep struct {
	repository.IRepository
}

// 储存下单信息，一笔订单对应一个下单信息
func (this *PaymentRep) Store(ctx context.Context, payment *models.Payment) (err error) {
	_, err = this.Collection().UpdateOne(ctx, bson.M{"order_no": payment.OrderNo}, bson.M{
		"$set": bson.M{
			"platform":         payment.Platform,
			"title":            payment.Title,
			"amount":           payment.Amount,
			"extended_user_id": payment.ExtendedUserId,
			"pre_payment_no":   payment.PrePaymentNo,
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

// 通过OrderNo获取支付信息
func (this *PaymentRep) FindByOrderId(ctx context.Context, orderNo string) (payment *models.Payment, err error) {
	result := this.Collection().FindOne(ctx, bson.M{"order_no": orderNo})
	if result.Err() != nil {
		return nil, err2.Err404.F("payment[order_no=%s],not found", orderNo)
	}
	payment = &models.Payment{}
	if err := result.Decode(payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (this *PaymentRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "order_no", Value: bsonx.Int64(-1)}},
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
	}
}

func NewPaymentRep(IRepository repository.IRepository) *PaymentRep {
	repository := &PaymentRep{IRepository: IRepository}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Printf("model [%s] create indexs error:%s\n", repository.TableName(), err)
		panic(err)
	}
	return repository
}
