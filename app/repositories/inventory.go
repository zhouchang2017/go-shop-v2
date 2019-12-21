package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
)

func init() {
	register(NewInventoryRep)
}

type InventoryRep struct {
	*mongoRep
}

func (this *InventoryRep) IncQty(ctx context.Context, filter interface{}, qty int64) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		update := this.Collection().FindOneAndUpdate(ctx,
			filter,
			bson.M{
				"$inc": bson.M{"qty": qty},
				"$currentDate": bson.M{
					"updated_at": true,
				},
			}, options.FindOneAndUpdate().SetReturnDocument(options.After))

		if update.Err() != nil {
			result <- repository.QueryResult{Error: update.Err()}
			return
		}
		model := this.newModel()
		err := update.Decode(model)
		result <- repository.QueryResult{Error: err, Result: model}
	}()
	return result
}

type QueryOption struct {
	ItemCode    string
	ItemId      string
	ProductId   string
	ProductCode string
	ShopId      string
	Location    *models.Location
	Status      *int8
}

func (q *QueryOption) SetStatus(status int8) {
	q.Status = &status
}

func (this *InventoryRep) Query(ctx context.Context, opt *QueryOption) {
	filter := bson.M{}
	if opt.ItemCode != "" {
		filter["item.code"] = opt.ItemCode
	}
	if opt.ItemId != "" {
		filter["item.id"] = opt.ItemId
		delete(filter, "item.code")
	}
	if opt.ProductCode != "" {
		filter["product.code"] = opt.ProductCode
	}
	if opt.ProductId != "" {
		filter["product.id"] = opt.ProductId
		delete(filter, "product.code")
	}
	if opt.ShopId != "" {
		filter["shop.id"] = opt.ShopId
	}
	// 状态过滤
	if opt.Status != nil {
		filter["status"] = *opt.Status
	}
	// 位置搜索
	if opt.Location != nil {
		filter["shop.location"] = bson.M{
			"$near":        opt.Location.GeoJSON(),
			"$maxDistance": 1000,
		}
	}
	this.Collection().Find(ctx, filter, options.Find().SetSort(bson.M{"qty": -1}))
}

// 索引
func (this *InventoryRep) indexesModel() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.M{
				"product.code":        "text",
				"product.name":        "text",
				"item.code":           "text",
				"shop.name":           "text",
				"option_values.value": "text",
			},
			Options: options.Index().SetWeights(bson.M{
				"product.code":        8,
				"product.name":        5,
				"item.code":           10,
				"shop.name":           5,
				"option_values.value": 3,
			}).SetBackground(true),
		},
		{
			Keys:    bsonx.Doc{{Key: "qty", Value: bsonx.Int64(-1)}},
			Options: options.Index().SetBackground(true),
		},
		{
			Keys:    bsonx.Doc{{Key: "shop.location", Value: bsonx.String("2d")}},
			Options: options.Index().SetBackground(true),
		},
	}
}

func NewInventoryRep(con *mongodb.Connection) *InventoryRep {
	rep := &InventoryRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Inventory{}, con),
	}
	err := rep.CreateIndexes(context.Background(), rep.indexesModel())
	if err != nil {
		log.Fatal(err)
	}
	return rep
}
