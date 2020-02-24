package repository

import (
	"context"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type InsertResult struct {
	Id     string
	Result interface{}
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
	FindById(ctx context.Context, id string) <-chan QueryResult
	FindByIds(ctx context.Context, ids ...string) <-chan QueryResult
	FindOne(ctx context.Context,credentials map[string]interface{}) <-chan QueryResult
	FindAll(ctx context.Context) <-chan QueryResult
	Count(ctx context.Context, filter interface{}) <-chan QueryResult
	Pagination(ctx context.Context, req *request.IndexRequest) <-chan QueryPaginationResult
	Create(ctx context.Context, entity interface{}) <-chan InsertResult
	Save(ctx context.Context, entity interface{}) <-chan QueryResult
	Update(ctx context.Context, id string, update interface{}) <-chan QueryResult
	Delete(ctx context.Context, id string) <-chan error
	DeleteMany(ctx context.Context, ids ...string) <-chan error
	Restore(ctx context.Context, id string) <-chan QueryResult
}
