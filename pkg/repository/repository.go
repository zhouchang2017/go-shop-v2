package repository

import (
	"context"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InsertResult struct {
	Id     string
	Result interface{}
	Error  error
}

type CountResult struct {
	Result int64
	Error  error
}

type QueryResult struct {
	Result interface{}
	Error  error
}

type QueryPaginationResult struct {
	Result     interface{}
	Error      error
	Pagination response.Pagination
}

type IRepository interface {
	TableName() string
	Collection() *mongo.Collection
	FindById(ctx context.Context, id string, opts ...*options.FindOneOptions) <-chan QueryResult
	FindByIds(ctx context.Context, ids []string, opts ...*options.FindOptions) <-chan QueryResult
	FindOne(ctx context.Context, credentials map[string]interface{}, opts ...*options.FindOneOptions) <-chan QueryResult
	FindAll(ctx context.Context) <-chan QueryResult
	FindMany(ctx context.Context, credentials map[string]interface{}, opts ...*options.FindOptions) <-chan QueryResult
	Count(ctx context.Context, filter interface{}) <-chan CountResult
	Pagination(ctx context.Context, req *request.IndexRequest) <-chan QueryPaginationResult
	AggregatePagination(ctx context.Context, entities interface{}, req *request.IndexRequest, pipe ...bson.D) <-chan QueryPaginationResult
	Create(ctx context.Context, entity interface{}) <-chan InsertResult
	Save(ctx context.Context, entity interface{}) <-chan QueryResult
	Update(ctx context.Context, id string, update interface{}) <-chan QueryResult
	Delete(ctx context.Context, id string) <-chan error
	DeleteMany(ctx context.Context, ids ...string) <-chan error
	Restore(ctx context.Context, id string) <-chan QueryResult
}
