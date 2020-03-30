package repositories

import (
	"context"
	"go-shop-v2/app/models"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
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

// 根据状态统计
func (this *OrderRep) CountByStatus(ctx context.Context, status int) (count int64, err error) {
	result := <-this.Count(ctx,
		bson.M{
			"status": status,
		})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.Result, nil
}

// 根据状态聚合统计数量
func (this *OrderRep) AggregateCountByStatus(ctx context.Context, statusList ...int) (response []*models.OrderCountByStatus, err error) {
	response = make([]*models.OrderCountByStatus, 0)
	if len(statusList) == 0 {
		return
	}
	aggregate, err := this.Collection().Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match", bson.M{"status": bson.M{"$in": statusList}}}},
		bson.D{{"$group", bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}}},
	})
	if err != nil {
		return response, err
	}
	if err := aggregate.All(ctx, &response); err != nil {
		return response, err
	}
	return
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

// 展开订单，订单中每件发货商品对应的门店信息
func (this *OrderRep) AggregateOrderItem(ctx context.Context, req *request.IndexRequest) (res []*models.AggregateOrderItem, pagination response.Pagination, err error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$unwind", "$logistics"}},
		bson.D{{"$unwind", "$logistics.items"}},
		bson.D{{"$unwind", "$order_items"}},
		bson.D{{"$replaceRoot", bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				bson.M{
					"order_id":       "$_id",
					"created_at":     "$created_at",
					"updated_at":     "$updated_at",
					"order_no":       "$order_no",
					"item_count":     "$item_count",
					"order_amount":   "$order_amount",
					"actual_amount":  "$actual_amount",
					"order_item":     "$order_items",
					"user":           "$user",
					"user_address":   "$user_address",
					"take_good_type": "$take_good_type",
					"logistics":      "$logistics",
					"payment":        "$payment",
					"status":         "$status",
					"promotion_info": "$promotion_info",
					"shipments_at":   "$shipments_at",
					"commented_at":   "$commented_at",
				},
			}},
		}}},
	}
	res = make([]*models.AggregateOrderItem, 1)
	result := <-this.IRepository.AggregatePagination(ctx, &res, req, pipeline...)
	if result.Error != nil {
		return nil, pagination, result.Error
	}
	data := result.Result.(*[]*models.AggregateOrderItem)
	res = *data
	pagination = result.Pagination
	return
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
