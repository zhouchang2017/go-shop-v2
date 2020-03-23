package repositories

import (
	"context"
	"go-shop-v2/app/models"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

type OrderRep struct {
	repository.IRepository
}

// 通过订单号查询订单
func (this *OrderRep) FindByOrderNo(ctx context.Context, orderNo string) (order *models.Order, err error) {
	result := this.Collection().FindOne(ctx, bson.M{"order_no": orderNo})
	if result.Err() != nil {
		return nil, err2.Err404.F("order[%s] not fond", orderNo)
	}
	order = &models.Order{}
	if err := result.Decode(order); err != nil {
		return nil, err
	}
	return
}

// 查询订单状态
func (this *OrderRep) GetOrderStatus(ctx context.Context, id string) (status int, err error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err2.Err404.F("order id[%s],not found", id)
	}
	one := this.Collection().FindOne(ctx, bson.M{"_id": objId}, options.FindOne().SetProjection(bson.M{
		"status": 1,
	}))
	if one.Err() != nil {
		return -1, err2.Err404.F("order id[%s],not found", id)
	}
	var order models.Order
	if err := one.Decode(&order); err != nil {
		return -1, err
	}
	return order.Status, nil
}

func (this *OrderRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "order_no", Value: bsonx.Int64(-1)}}, // order_no 唯一键
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
	}
}

func NewOrderRep(rep repository.IRepository) *OrderRep {
	repository := &OrderRep{rep}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Printf("model [%s] create indexs error:%s\n", repository.TableName(), err)
		panic(err)
	}
	return repository
}
