package repositories

import (
	"context"
	"errors"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/utils"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type mongoRep struct {
	table string
	model interface{}
	con   *mongodb.Connection
}

func NewBasicMongoRepositoryByDefault(model interface{}, con *mongodb.Connection) *mongoRep {
	plural := utils.StructNameToSnakeAndPlural(model)
	return &mongoRep{table: plural, model: model, con: con}
}

func (this mongoRep) TableName() string {
	return this.table
}

func (this *mongoRep) newModel() interface{} {
	t := reflect.ValueOf(this.model).Elem().Type()
	return reflect.New(t).Interface()
}

func (this *mongoRep) newModels() interface{} {
	modelType := reflect.TypeOf(this.model)
	slice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0)
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)
	return x.Interface()
	//return reflect.New(reflect.SliceOf(t)).Interface()
}

func (this *mongoRep) Collection() *mongo.Collection {
	return this.con.Collection(this.table)
}

func (this *mongoRep) FindMany(ctx context.Context, credentials map[string]interface{}) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	trashed := ctx2.GetTrashed(ctx)
	go func() {
		defer close(output)
		filter := bson.M{}
		if !trashed {
			filter["deleted_at"] = bson.D{{"$eq", nil}}
		}

		for k, v := range credentials {
			filter[k] = v
		}
		cursor, err := this.Collection().Find(ctx, filter)
		if err != nil {
			output <- repository.QueryResult{Error: err}
			return
		}
		entities := this.newModels()
		err = cursor.All(ctx, entities)
		if err != nil {
			output <- repository.QueryResult{Error: err}
			return
		}

		i := reflect.ValueOf(entities).Elem().Interface()

		output <- repository.QueryResult{Result: i}
	}()
	return output
}
func (this *mongoRep) FindOne(ctx context.Context, credentials map[string]interface{}) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	trashed := ctx2.GetTrashed(ctx)
	go func() {
		defer close(output)
		filter := bson.M{}
		if !trashed {
			filter["deleted_at"] = bson.D{{"$eq", nil}}
		}

		for k, v := range credentials {
			filter[k] = v
		}
		one := this.Collection().FindOne(ctx, filter)
		if one.Err() != nil {
			if one.Err() == mongo.ErrNoDocuments {
				output <- repository.QueryResult{Error: err2.Err404}
				return
			}
			output <- repository.QueryResult{Error: one.Err()}
			return
		}
		entity := this.newModel()
		err := one.Decode(entity)
		if err != nil {
			output <- repository.QueryResult{Error: err}
			return
		}
		output <- repository.QueryResult{Result: entity}
	}()
	return output
}

func (this *mongoRep) FindById(ctx context.Context, id string) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	trashed := ctx2.GetTrashed(ctx)
	go func() {
		defer close(output)

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			output <- repository.QueryResult{Error: err2.Err404}
			return
		}
		filter := bson.M{"_id": objId}
		if !trashed {
			filter["deleted_at"] = bson.D{{"$eq", nil}}
		}

		one := this.Collection().FindOne(ctx, filter)
		if one.Err() != nil {
			if one.Err() == mongo.ErrNoDocuments {
				output <- repository.QueryResult{Error: err2.Err404}
				return
			}
			output <- repository.QueryResult{Error: one.Err()}
			return
		}
		entity := this.newModel()
		err = one.Decode(entity)
		if err != nil {
			output <- repository.QueryResult{Error: err}
			return
		}

		output <- repository.QueryResult{Result: entity}
	}()
	return output
}

func (this *mongoRep) FindByIds(ctx context.Context, ids ...string) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	trashed := ctx2.GetTrashed(ctx)
	go func() {
		defer close(output)
		models := this.newModels()
		var objIds []primitive.ObjectID
		for _, id := range ids {
			objectIDS, e := primitive.ObjectIDFromHex(id)
			if e == nil {
				objIds = append(objIds, objectIDS)
			}
		}
		if len(objIds) == 0 {
			output <- repository.QueryResult{Result: reflect.ValueOf(models).Elem().Interface()}
			return
		}
		filter := bson.M{"_id": bson.M{"$in": objIds}}
		if !trashed {
			filter["deleted_at"] = bson.D{{"$eq", nil}}
		}

		cursor, err := this.Collection().Find(ctx, filter)
		if err != nil {
			output <- repository.QueryResult{Error: err}
			return
		}

		err = cursor.All(ctx, models)
		if err != nil {
			output <- repository.QueryResult{Error: err}
			return
		}
		i := reflect.ValueOf(models).Elem().Interface()
		output <- repository.QueryResult{Result: i}
	}()
	return output
}

func (this *mongoRep) FindAll(ctx context.Context) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	trashed := ctx2.GetTrashed(ctx)

	go func() {
		defer close(output)

		filter := bson.M{}
		if !trashed {
			filter["deleted_at"] = bson.D{{"$eq", nil}}
		}

		cursor, err := this.Collection().Find(ctx, filter)
		if err != nil {
			output <- repository.QueryResult{Error: err}
			return
		}
		entities := this.newModels()
		err = cursor.All(ctx, entities)
		if err != nil {
			output <- repository.QueryResult{Error: err}
			return
		}

		i := reflect.ValueOf(entities).Elem().Interface()

		output <- repository.QueryResult{Result: i}
	}()
	return output
}

func (this *mongoRep) Count(ctx context.Context, filter interface{}) <-chan repository.CountResult {
	output := make(chan repository.CountResult)
	go func() {
		defer close(output)

		if filter == nil {
			filter = bson.M{}
		}

		total, err := this.Collection().CountDocuments(ctx, filter)
		if err != nil {
			output <- repository.CountResult{Error: err}
			return
		}

		output <- repository.CountResult{Result: total}
	}()
	return output
}

func (this *mongoRep) Pagination(ctx context.Context, req *request.IndexRequest) <-chan repository.QueryPaginationResult {
	result := make(chan repository.QueryPaginationResult)

	go func() {
		defer close(result)
		page := req.GetPage()
		find := options.Find()
		if page > 0 {
			// 分页
			find.SetLimit(req.GetPerPage())
			find.SetSkip((page - 1) * req.GetPerPage())
		}
		// 排序
		if opt, ok := req.Sort(); ok {
			find.SetSort(opt)
		} else {
			// 默认id降序
			find.SetSort(bson.M{"_id": -1})
		}
		// SELECT
		if opt, ok := req.Projection(); ok {
			find.SetProjection(opt)
		}
		// query builder
		filter := bson.M{}
		// search
		if req.Search != "" && req.GetSearchField() != "" {
			filter[req.GetSearchField()] = primitive.Regex{Pattern: req.Search, Options: "i"}
		}
		// trashed
		if !req.Trashed {
			filter["deleted_at"] = bson.D{{"$eq", nil}}
		}
		// 自定义过滤
		for key, value := range req.Query() {
			filter[key] = value
		}
		// filter ids
		if req.Ids != nil {
			var objIds []primitive.ObjectID
			for _, id := range strings.Split(*req.Ids, ",") {
				if ids, err := primitive.ObjectIDFromHex(id); err == nil {
					objIds = append(objIds, ids)
				}
			}
			if len(objIds) > 0 {
				filter["_id"] = bson.M{"$in": objIds}
			} else {
				filter["_id"] = bson.M{"$in": bson.A{}}
			}
		}

		var total int64
		// count
		r := <-this.Count(ctx, filter)
		if r.Error != nil {
			result <- repository.QueryPaginationResult{Error: r.Error}
			return
		}
		total = r.Result

		cursor, err := this.Collection().Find(ctx, filter, find)
		if err != nil {
			result <- repository.QueryPaginationResult{Error: err}
			return
		}

		entities := this.newModels()
		err = cursor.All(ctx, entities)
		if err != nil {
			result <- repository.QueryPaginationResult{Error: err}
			return
		}
		pagination := response.Pagination{
			Total: total,
		}
		if req.Page != -1 {
			pagination.CurrentPage = page
			pagination.PerPage = req.PerPage
			pagination.HasNextPage = page*req.PerPage < total
		}

		i := reflect.ValueOf(entities).Elem().Interface()

		result <- repository.QueryPaginationResult{Result: i, Pagination: pagination}
	}()
	return result
}

func (this *mongoRep) Create(ctx context.Context, entity interface{}) <-chan repository.InsertResult {
	result := make(chan repository.InsertResult)
	go func() {
		defer close(result)
		this.setId(entity)
		this.setTimeNow(entity, "CreatedAt")
		this.setTimeNow(entity, "UpdatedAt")

		oneResult, err := this.Collection().InsertOne(ctx, entity)
		if err != nil {
			result <- repository.InsertResult{Error: err}
			return
		}
		insertResult := repository.InsertResult{}
		if obj, ok := oneResult.InsertedID.(primitive.ObjectID); ok {
			this.setValue(entity, "ID", obj)
			insertResult.Id = obj.Hex()
			findResult := <-this.FindById(ctx, obj.Hex())
			if findResult.Error != nil {
				result <- repository.InsertResult{Error: findResult.Error}
				return
			}
			insertResult.Result = findResult.Result
		}
		result <- insertResult
	}()
	return result
}

func (this *mongoRep) Save(ctx context.Context, entity interface{}) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		ids, err := this.getId(entity)
		if err != nil {
			result <- repository.QueryResult{Error: err2.Err404}
			return
		}
		this.setTimeNow(entity, "UpdatedAt")
		update := this.Collection().FindOneAndUpdate(ctx,
			bson.M{"_id": ids},
			bson.M{
				"$set": entity,
			}, options.FindOneAndUpdate().SetReturnDocument(options.After))

		if update.Err() != nil {
			result <- repository.QueryResult{Error: update.Err()}
			return
		}
		model := this.newModel()
		err = update.Decode(model)
		if err != nil {
			result <- repository.QueryResult{Error: err}
			return
		}

		result <- repository.QueryResult{Result: model}
	}()
	return result
}

func (this *mongoRep) Update(ctx context.Context, id string, update interface{}) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		objid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			result <- repository.QueryResult{Error: err2.Err404}
			return
		}

		updated := bson.M{}

		of := reflect.ValueOf(update)
		switch of.Kind() {
		case reflect.Map:
			for _, key := range of.MapKeys() {
				updated[key.String()] = of.MapIndex(key).Interface()
			}
		}

		updated["$currentDate"] = bson.M{
			"updated_at": true,
		}

		update := this.Collection().FindOneAndUpdate(ctx,
			bson.M{"_id": objid},
			updated, options.FindOneAndUpdate().SetReturnDocument(options.After))

		if update.Err() != nil {
			result <- repository.QueryResult{Error: update.Err()}
			return
		}
		model := this.newModel()
		err = update.Decode(model)
		if err != nil {
			result <- repository.QueryResult{Error: err}
			return
		}

		result <- repository.QueryResult{Result: model}
	}()
	return result
}

func (this *mongoRep) destroy(ctx context.Context, id primitive.ObjectID) <-chan error {
	result := make(chan error)
	go func() {
		defer close(result)
		_, err := this.Collection().DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			result <- err
			return
		}
		result <- nil
		return
	}()
	return result
}

func (this *mongoRep) Delete(ctx context.Context, id string) <-chan error {
	result := make(chan error)
	go func() {
		defer close(result)

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			result <- err2.Err404
			return
		}
		force := ctx2.GetForce(ctx)
		if force {
			// 硬删除
			destroy := <-this.destroy(ctx, objId)
			if destroy != nil {
				result <- destroy
				return
			}
			result <- nil
			return
		}
		if this.hasDeletedAtField(this.model) {
			// soft delete
			now := time.Now()
			this.setValue(this.model, "DeletedAt", &now)
			_, err := this.Collection().UpdateOne(ctx, bson.M{"_id": objId}, bson.M{
				"$set": bson.M{"deleted_at": now},
				"$currentDate": bson.M{
					"updated_at": true,
				},
			})
			if err != nil {
				result <- err
				return
			}
			result <- nil
			return
		} else {
			destroy := <-this.destroy(ctx, objId)
			if destroy != nil {
				result <- destroy
				return
			}
			result <- nil
			return
		}
	}()
	return result
}

func (this *mongoRep) DeleteMany(ctx context.Context, ids ...string) <-chan error {
	result := make(chan error)
	go func() {
		defer close(result)

		var objIds []primitive.ObjectID
		for _, id := range ids {
			objId, err := primitive.ObjectIDFromHex(id)
			if err == nil {
				objIds = append(objIds, objId)

			}
		}
		if len(objIds) == 0 {
			result <- nil
			return
		}

		force := ctx2.GetForce(ctx)

		if this.hasDeletedAtField(this.model) && !force {
			// soft delete
			now := time.Now()
			//this.setValue(this.model, "DeletedAt", &now)
			_, err := this.Collection().UpdateMany(ctx, bson.M{"_id": bson.M{"$in": objIds}}, bson.M{
				"$set": bson.M{"deleted_at": now},
				"$currentDate": bson.M{
					"updated_at": true,
				},
			})
			if err != nil {
				result <- err
				return
			}

			result <- nil
			return
		} else {
			// 硬删除
			_, err := this.Collection().DeleteMany(ctx, bson.M{
				"_id": bson.M{"$in": objIds},
			})

			if err != nil {
				result <- err
				return
			}

			result <- nil
			return
		}
	}()
	return result
}

func (this *mongoRep) Restore(ctx context.Context, id string) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			result <- repository.QueryResult{Error: err2.Err404}
			return
		}
		update := this.Collection().FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{
			"$set": bson.M{"deleted_at": nil},
			"$currentDate": bson.M{
				"updated_at": true,
			},
		}, options.FindOneAndUpdate().SetReturnDocument(options.After))
		if update.Err() != nil {
			result <- repository.QueryResult{Error: update.Err()}
			return
		}
		model := this.newModel()
		err = update.Decode(model)
		result <- repository.QueryResult{Error: err, Result: model}
	}()
	return result
}

type AggregateMetadata struct {
	Total int64 `bson:"total"`
}

// 聚合分页
func (this *mongoRep) AggregatePagination(ctx context.Context, entities interface{}, req *request.IndexRequest, pipe ...bson.D) <-chan repository.QueryPaginationResult {
	result := make(chan repository.QueryPaginationResult)

	go func() {
		defer close(result)
		// 初始化
		pipeline := mongo.Pipeline{}

		// match阶段
		// query builder
		filter := bson.M{}
		// trashed
		if !req.Trashed {
			filter["deleted_at"] = bson.D{{"$eq", nil}}
		}
		// search
		if req.Search != "" && req.GetSearchField() != "" {
			filter[req.GetSearchField()] = primitive.Regex{Pattern: req.Search, Options: "i"}
		}

		// 自定义过滤
		for key, value := range req.Query() {
			filter[key] = value
		}

		if len(filter) > 0 {
			pipeline = append(pipeline, bson.D{{"$match", filter}})
		}

		// 合并 groupBy 阶段
		pipeline = append(pipeline, pipe...)

		// SELECT
		if opt, ok := req.Projection(); ok {
			pipeline = append(pipeline, bson.D{{"$project", opt}})
		}

		// 排序
		if opt, ok := req.Sort(); ok {
			pipeline = append(pipeline, bson.D{{"$sort", opt}})
		}

		// 当前页
		page := req.GetPage()
		// facet
		facet := bson.M{
			"metadata": bson.A{bson.D{{"$count", "total"}}},
		}
		limit := req.GetPerPage()
		skip := (page - 1) * req.GetPerPage()
		if page > 0 {
			// 分页
			facet["data"] = bson.A{
				bson.D{{"$skip", skip}},
				bson.D{{"$limit", limit}},
			}
		} else {
			//data: [{ $replaceRoot: { newRoot: "$$ROOT" } }]
			facet["data"] = bson.A{
				bson.D{{"$replaceRoot", bson.D{{"newRoot", "$$ROOT"}}}},
			}
		}

		pipeline = append(pipeline, bson.D{{"$facet", facet}})

		cursor, err := this.Collection().Aggregate(ctx, pipeline)
		if err != nil {
			result <- repository.QueryPaginationResult{Error: err}
			return
		}
		var total int64 = 0
		var metadata []AggregateMetadata
		defer cursor.Close(ctx)

		if cursor.Next(ctx) {
			lookup := cursor.Current.Lookup("data")
			err := lookup.Unmarshal(entities)
			if err != nil {
				result <- repository.QueryPaginationResult{Error: err}
				return
			}

			value := cursor.Current.Lookup("metadata")
			err = value.Unmarshal(&metadata)
			if err != nil {
				result <- repository.QueryPaginationResult{Error: err}
				return
			}
			if len(metadata) > 0 {
				total = metadata[0].Total
			}
		}

		pagination := response.Pagination{
			Total: total,
		}
		if page != -1 {
			pagination.CurrentPage = page
			pagination.PerPage = limit
			pagination.HasNextPage = page*limit < total
		}

		result <- repository.QueryPaginationResult{Result: entities, Pagination: pagination}
	}()
	return result
}

func (this *mongoRep) setTimeNow(entity interface{}, field string) {
	this.setValue(entity, field, time.Now())
}

func (this *mongoRep) setId(entity interface{}) {
	if reflect.ValueOf(entity).Kind() == reflect.Ptr {
		elem := reflect.ValueOf(entity).Elem()
		f := elem.FieldByName("ID")
		if f.IsValid() {
			if f.Type() == reflect.ValueOf(primitive.ObjectID{}).Type() {
				if f.Interface().(primitive.ObjectID).IsZero() {
					f.Set(reflect.ValueOf(primitive.NewObjectID()))
				}
			}
		}
	}
}

func (this *mongoRep) getId(entity interface{}) (primitive.ObjectID, error) {
	if reflect.ValueOf(entity).Kind() == reflect.Ptr {
		elem := reflect.ValueOf(entity).Elem()
		f := elem.FieldByName("ID")
		if f.IsValid() {
			if f.Type() == reflect.ValueOf(primitive.ObjectID{}).Type() {
				if !f.Interface().(primitive.ObjectID).IsZero() {
					return f.Interface().(primitive.ObjectID), nil
				}
			}
		}
	}
	return primitive.ObjectID{}, errors.New("can not find id")
}

func (this *mongoRep) setValue(entity interface{}, field string, value interface{}) {
	if reflect.ValueOf(entity).Kind() == reflect.Ptr {
		elem := reflect.ValueOf(entity).Elem()
		f := elem.FieldByName(field)
		if f.IsValid() {
			if f.Type() == reflect.ValueOf(value).Type() {
				f.Set(reflect.ValueOf(value))
			}
		}
	}
}

func (this *mongoRep) hasDeletedAtField(entity interface{}) bool {
	if reflect.ValueOf(entity).Kind() == reflect.Ptr {
		elem := reflect.ValueOf(entity).Elem()
		f := elem.FieldByName("DeletedAt")
		if f.IsValid() {
			if f.Type() == reflect.ValueOf(&time.Time{}).Type() {
				return true
			}
		}
	}
	return false
}

func (this *mongoRep) CreateIndexes(ctx context.Context, models []mongo.IndexModel) (err error) {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err = this.Collection().Indexes().CreateMany(ctx, models, opts)
	if err != nil {
		log.Printf("model %s create indexs error:%s\n", this.table, err)
		return err
	}
	//for _, key := range res {
	//	log.Printf("model %s create indexs %s\n", this.table, key)
	//}

	return nil
}
