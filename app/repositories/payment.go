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

// 获取某天收款总额
func (this *PaymentRep) GetRangePaymentCount(ctx context.Context, start time.Time, end time.Time) (response *models.DayPaymentCount, err error) {
	response = &models.DayPaymentCount{
		TotalAmount: 0,
		Count:       0,
	}
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.M{
			"payment_at": bson.M{"$gte": start, "$lte": end},
		}}},
		bson.D{{"$group",
			bson.M{
				"_id":          nil,
				"total_amount": bson.M{"$sum": "$amount"},
				"count":        bson.M{"$sum": 1},
			},
		}},
		bson.D{{"$replaceRoot", bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				bson.M{
					"total_amount": "$total_amount",
					"count":        "$count",
				},
			}},
		}}},
	}
	aggregate, err := this.Collection().Aggregate(ctx, pipeline)
	if err != nil {
		return response, err
	}
	var res []*models.DayPaymentCount
	if err := aggregate.All(ctx, &res); err != nil {
		return response, err
	}
	if len(res) == 1 {
		return res[0], nil
	}
	return response, nil
}

// 获取一段时间收款统计
func (this *PaymentRep) GetRangePaymentCounts(ctx context.Context, start time.Time, end time.Time) (response []*models.DayPaymentCount, err error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.M{
			"payment_at": bson.M{"$gte": start, "$lte": end},
		}}},
		bson.D{{"$group",
			bson.M{
				"_id": bson.M{"$dateToString": bson.M{
					"format": "%Y-%m-%d",
					"date":   "$payment_at",
				}},
				"total_amount": bson.M{"$sum": "$amount"},
				"count":        bson.M{"$sum": 1},
			},
		}},
		bson.D{{"$sort", bson.M{"_id": 1}}},
	}
	aggregate, err := this.Collection().Aggregate(ctx, pipeline)
	if err != nil {
		return response, err
	}
	response = make([]*models.DayPaymentCount, 0)
	if err := aggregate.All(ctx, &response); err != nil {
		return response, err
	}
	return response, nil
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
