package repositories

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func init() {
	register(NewInventoryRep)
}

type InventoryRep struct {
	*mongoRep
}

func (this *InventoryRep) AggregateStockByShops(ctx context.Context, shopIds ...string) (data []*models.AggregateShopCountStockInventory, err error) {

	var objIds []primitive.ObjectID
	for _, id := range shopIds {
		if ids, err := primitive.ObjectIDFromHex(id); err == nil {
			objIds = append(objIds, ids)
		}
	}

	pipelines1 := mongo.Pipeline{}

	if len(objIds) > 0 {
		pipelines1 = append(pipelines1, bson.D{{"$match", bson.M{
			"shop.id": bson.D{{"$in", objIds}},
		}}})
	}

	pipelines2 := mongo.Pipeline{
		// 状态分组统计
		bson.D{{"$group", bson.M{
			"_id": bson.M{
				"shop_id":   "$shop.id",
				"shop_name": "$shop.name",
				"status":    "$status",
			},
			"qty": bson.D{{"$sum", "$qty"}},
		}}},
		// 合并状态统计
		bson.D{{"$group", bson.M{
			"_id": bson.M{
				"shop_id":   "$_id.shop_id",
				"shop_name": "$_id.shop_name",
			},
			"total": bson.D{{"$sum", "$qty"}},
			"status": bson.D{{"$push", bson.M{
				"status": "$_id.status",
				"qty":    "$qty",
			}}},
		}}},
		// 改变数据结构
		bson.D{{"$replaceRoot",
			bson.D{{"newRoot",
				bson.D{{"$mergeObjects",
					bson.A{bson.M{
						"shop_id":   "$_id.shop_id",
						"shop_name": "$_id.shop_name",
						"total":     "$total",
						"status":    "$status",
					}}}}}}}},
	}

	pipelines1 = append(pipelines1, pipelines2...)

	cursor, err := this.Collection().Aggregate(ctx, pipelines1)
	if err != nil {
		return nil, err
	}

	data = []*models.AggregateShopCountStockInventory{}
	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (this *InventoryRep) AggregatePagination(ctx context.Context, req *request.IndexRequest) <-chan repository.QueryPaginationResult {
	result := make(chan repository.QueryPaginationResult)

	go func() {
		defer close(result)

		es := []*models.AggregateInventory{}
		filters := req.Filters.Unmarshal()
		options := &QueryOption{}
		err := mapstructure.Decode(filters, options)

		if err != nil {
			result <- repository.QueryPaginationResult{Error: err}
			return
		}
		if len(options.Status) > 0 {
			req.AppendFilter("status", bson.D{{"$in", options.Status}})
		}
		if len(options.Shops) > 0 {
			req.AppendFilter("shop.id", bson.D{{"$in", options.Shops}})
		}

		pipelines := mongo.Pipeline{
			// 状态分组统计
			bson.D{{"$group", bson.M{
				"_id": bson.M{
					"item_id": "$item.id",
					"status":  "$status",
				},
				"item": bson.D{{"$mergeObjects", "$item"}},
				"qty":  bson.D{{"$sum", "$qty"}},
				"shops": bson.D{{"$push", bson.M{
					"id":           "$shop.id",
					"name":         "$shop.name",
					"qty":          "$qty",
					"inventory_id": "$_id",
				}}},
			}}},
			// 合并状态统计
			bson.D{{"$group", bson.M{
				"_id": bson.M{
					"item_id": "$_id.item_id",
				},
				"item":  bson.D{{"$mergeObjects", "$item"}},
				"total": bson.D{{"$sum", "$qty"}},
				"inventories": bson.D{{"$push", bson.M{
					"status": "$_id.status",
					"qty":    "$qty",
					"shops":  "$shops",
				}}},
			}}},
			// 改变数据结构
			bson.D{{"$replaceRoot",
				bson.D{{"newRoot",
					bson.D{{"$mergeObjects",
						bson.A{"$item", bson.M{
							"total":       "$total",
							"inventories": "$inventories",
						}}}}}}}},
		}
		aggregateRes := <-this.mongoRep.AggregatePagination(ctx, &es, req, pipelines...)
		if aggregateRes.Error != nil {
			result <- repository.QueryPaginationResult{Error: aggregateRes.Error}
			return
		}
		result <- repository.QueryPaginationResult{Result: es, Pagination: aggregateRes.Pagination, Error: aggregateRes.Error}
	}()

	return result
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
	Brands      []string `json:"brands"`
	ItemCode    string   `json:"item_code"`
	ItemId      string   `json:"item_id"`
	ProductId   string   `json:"product_id"`
	ProductCode string   `json:"product_code"`
	Shops       []string `json:"shops"`
	Location    *models.Location
	Status      []int8
}

//func (q *QueryOption) SetStatus(status int8) {
//	q.Status = &status
//}

func (this *InventoryRep) Pagination(ctx context.Context, req *request.IndexRequest) <-chan repository.QueryPaginationResult {
	result := make(chan repository.QueryPaginationResult)
	req.SetSearchField("item.code")
	filters := req.Filters.Unmarshal()
	options := &QueryOption{}
	err := mapstructure.Decode(filters, options)
	spew.Dump(options)
	if err != nil {
		defer close(result)
		result <- repository.QueryPaginationResult{Error: err}
		return result
	}
	if len(options.Status) > 0 {
		req.AppendFilter("status", bson.D{{"$in", options.Status}})
	}
	if len(options.Shops) > 0 {
		req.AppendFilter("shop.id", bson.D{{"$in", options.Shops}})
	}

	return this.mongoRep.Pagination(ctx, req)

	//if opt.ItemCode != "" {
	//	filter["item.code"] = opt.ItemCode
	//}
	//if opt.ItemId != "" {
	//	filter["item.id"] = opt.ItemId
	//	delete(filter, "item.code")
	//}
	//if opt.ProductCode != "" {
	//	filter["product.code"] = opt.ProductCode
	//}
	//if opt.ProductId != "" {
	//	filter["product.id"] = opt.ProductId
	//	delete(filter, "product.code")
	//}
	//if opt.ShopId != "" {
	//	filter["shop.id"] = opt.ShopId
	//}
	//// 状态过滤
	//if opt.Status != nil {
	//	filter["status"] = *opt.Status
	//}
	//// 位置搜索
	//if opt.Location != nil {
	//	filter["shop.location"] = bson.M{
	//		"$near":        opt.Location.GeoJSON(),
	//		"$maxDistance": 1000,
	//	}
	//}
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
	//err := rep.CreateIndexes(context.Background(), rep.indexesModel())
	//if err != nil {
	//	log.Fatal(err)
	//}
	return rep
}
