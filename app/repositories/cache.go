package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/cache/redis"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"reflect"
	"time"
)

type RedisCache struct {
	model interface{}
	*redis.Connection
	rep *mongoRep
}

func (this RedisCache) FindMany(ctx context.Context, credentials map[string]interface{}) <-chan repository.QueryResult {
	return this.rep.FindMany(ctx, credentials)
}

func (this RedisCache) AggregatePagination(ctx context.Context, entities interface{}, req *request.IndexRequest, pipe ...bson.D) <-chan repository.QueryPaginationResult {
	return this.rep.AggregatePagination(ctx, entities, req, pipe...)
}

func NewRedisCache(model interface{}, con *redis.Connection, rep *mongoRep) *RedisCache {
	return &RedisCache{model: model, Connection: con, rep: rep}
}

func (this RedisCache) Collection() *mongo.Collection {
	return this.rep.Collection()
}

func (this RedisCache) TableName() string {
	return this.rep.TableName()
}

func (this RedisCache) GetCacheKey(id string) string {
	return fmt.Sprintf("%s:%s", this.TableName(), id)
}

func (this *RedisCache) newModel() interface{} {
	t := reflect.ValueOf(this.model).Elem().Type()
	return reflect.New(t).Interface()
}

func (this *RedisCache) newModels() interface{} {
	t := reflect.ValueOf(this.model).Type()
	return reflect.New(reflect.SliceOf(t)).Interface()
}

func (this *RedisCache) getId(entity interface{}) (primitive.ObjectID, error) {
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

func (this *RedisCache) deletedAtFieldIsNil(entity interface{}) bool {
	if reflect.ValueOf(entity).Kind() == reflect.Ptr {
		elem := reflect.ValueOf(entity).Elem()
		f := elem.FieldByName("DeletedAt")
		if f.IsValid() {
			if f.Type() == reflect.ValueOf(&time.Time{}).Type() {
				if f.IsNil() {
					return true
				}
				return false
			}
		}
	}
	return true
}

func (this RedisCache) FindById(ctx context.Context, id string) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	trashed := ctx2.GetTrashed(ctx)
	go func() {
		defer close(output)
		result := this.HMGet(this.GetCacheKey(id), "detail").Val()
		if len(result) > 0 {
			if result[0] != nil {
				jsonValue := result[0].(string)
				// Unmarshal
				entity := this.newModel()
				if err := json.Unmarshal([]byte(jsonValue), entity); err != nil {
					log.Printf("%s [%s] FindById form cache,error:%s\n", this.TableName(), id, err)
					// 从缓存中移除
					this.HDel(this.GetCacheKey(id), "detail")
				}
				if !trashed {
					if !this.deletedAtFieldIsNil(entity) {
						output <- repository.QueryResult{Error: err2.Err404}
						return
					}
				}
				// 命中缓存
				output <- repository.QueryResult{Result: entity}
				return
			}
		}
		results := <-this.rep.FindById(ctx, id)
		if results.Error == nil {
			if marshal, err := json.Marshal(results.Result); err == nil {
				if _, err := this.HMSet(this.GetCacheKey(id), "detail", marshal).Result(); err != nil {
					spew.Dump(err)
				}
			}

		}
		output <- repository.QueryResult{
			Result: results.Result,
			Error:  results.Error,
		}
	}()
	return output
}

func (this RedisCache) FindByIds(ctx context.Context, ids ...string) <-chan repository.QueryResult {
	return this.rep.FindByIds(ctx, ids...)
}

func (this RedisCache) FindOne(ctx context.Context, credentials map[string]interface{}) <-chan repository.QueryResult {
	return this.rep.FindOne(ctx, credentials)
}

func (this RedisCache) FindAll(ctx context.Context) <-chan repository.QueryResult {
	return this.rep.FindAll(ctx)
}

func (this RedisCache) Count(ctx context.Context, filter interface{}) <-chan repository.CountResult {
	return this.rep.Count(ctx, filter)
}

func (this RedisCache) Pagination(ctx context.Context, req *request.IndexRequest) <-chan repository.QueryPaginationResult {
	return this.rep.Pagination(ctx, req)
}

func (this RedisCache) Create(ctx context.Context, entity interface{}) <-chan repository.InsertResult {
	result := make(chan repository.InsertResult)
	go func() {
		defer close(result)
		results := <-this.rep.Create(ctx, entity)
		if results.Error == nil {
			if marshal, err := json.Marshal(results.Result); err == nil {
				if _, err := this.HMSet(this.GetCacheKey(results.Id), "detail", marshal).Result(); err != nil {
					spew.Dump(err)
				}
			}

		}
		result <- repository.InsertResult{
			Id:     results.Id,
			Result: results.Result,
			Error:  results.Error,
		}
	}()

	return result
}

func (this RedisCache) Save(ctx context.Context, entity interface{}) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		results := <-this.rep.Save(ctx, entity)
		if results.Error == nil {
			if id, err := this.getId(results.Result); err == nil {
				if marshal, err := json.Marshal(results.Result); err == nil {
					if _, err := this.HMSet(this.GetCacheKey(id.Hex()), "detail", marshal).Result(); err != nil {
						spew.Dump(err)
					}
				}
			}
		}
		result <- repository.QueryResult{
			Result: results.Result,
			Error:  results.Error,
		}
	}()
	return result
}

func (this RedisCache) Update(ctx context.Context, id string, update interface{}) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		results := <-this.rep.Update(ctx, id, update)
		if results.Error == nil {
			if marshal, err := json.Marshal(results.Result); err == nil {
				if _, err := this.HMSet(this.GetCacheKey(id), "detail", marshal).Result(); err != nil {
					spew.Dump(err)
				}
			}
		}
		result <- repository.QueryResult{
			Result: results.Result,
			Error:  results.Error,
		}
	}()
	return result
}

func (this RedisCache) Delete(ctx context.Context, id string) <-chan error {
	result := make(chan error)
	go func() {
		defer close(result)
		err := <-this.rep.Delete(ctx, id)
		if err == nil {
			this.HDel(this.GetCacheKey(id), "detail")
		}
		result <- err
	}()
	return result
}

func (this RedisCache) DeleteMany(ctx context.Context, ids ...string) <-chan error {
	result := make(chan error)
	go func() {
		defer close(result)
		err := <-this.rep.DeleteMany(ctx, ids...)
		if err == nil {
			for _, id := range ids {
				this.HDel(this.GetCacheKey(id), "detail")
			}
		}
		result <- err
	}()
	return result
}

func (this RedisCache) Restore(ctx context.Context, id string) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		results := <-this.rep.Restore(ctx, id)
		if results.Error == nil {
			if marshal, err := json.Marshal(results.Result); err == nil {
				if _, err := this.HMSet(this.GetCacheKey(id), "detail", marshal).Result(); err != nil {
					spew.Dump(err)
				}
			}
		}
		result <- repository.QueryResult{
			Result: results.Result,
			Error:  results.Error,
		}
	}()
	return result
}
