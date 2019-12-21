package repositories

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"reflect"
)

type mysqlRep struct {
	model interface{}
	con   *gorm.DB
}

func (this *mysqlRep) FindOne(ctx context.Context, credentials map[string]interface{}) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	go func() {
		defer close(output)
		model := this.newModel()
		first := this.con.Where(credentials).First(model)
		output <- repository.QueryResult{Result: model, Error: first.Error}
	}()
	return output
}

func NewMysqlRep(model interface{}, con *gorm.DB) *mysqlRep {
	return &mysqlRep{model: model, con: con}
}

func (this *mysqlRep) newModel() interface{} {
	t := reflect.ValueOf(this.model).Elem().Type()
	return reflect.New(t).Interface()
}

func (this *mysqlRep) newModels() interface{} {
	t := reflect.ValueOf(this.model).Type()
	return reflect.New(reflect.SliceOf(t)).Interface()
}

func (this *mysqlRep) FindById(ctx context.Context, id string) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	go func() {
		defer close(output)
		model := this.newModel()
		first := this.con.First(model, id)
		output <- repository.QueryResult{Result: model, Error: first.Error}
	}()
	return output

}

func (this *mysqlRep) FindByIds(ctx context.Context, ids ...string) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	go func() {
		defer close(output)
		models := this.newModels()
		find := this.con.Where("id in (?)", ids).Find(models)
		i := reflect.ValueOf(models).Elem().Interface()
		output <- repository.QueryResult{Result: i, Error: find.Error}
	}()
	return output
}

func (this *mysqlRep) FindAll(ctx context.Context) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	go func() {
		defer close(output)
		models := this.newModels()
		find := this.con.Find(models)
		i := reflect.ValueOf(models).Elem().Interface()
		output <- repository.QueryResult{Result: i, Error: find.Error}
	}()
	return output
}

func (this *mysqlRep) Count(ctx context.Context, filter interface{}) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	go func() {
		defer close(output)
		var total int64
		var res *gorm.DB
		if filter != nil {
			res = this.con.Model(this.model).Where(filter).Count(&total)
		} else {
			res = this.con.Model(this.model).Count(&total)
		}

		output <- repository.QueryResult{Result: total, Error: res.Error}
	}()
	return output
}

func (this *mysqlRep) Pagination(ctx context.Context, req *request.IndexRequest) <-chan repository.QueryPaginationResult {
	result := make(chan repository.QueryPaginationResult)
	go func() {
		defer close(result)
		db := this.con

		// soft delete
		if req.Trashed {
			db = db.Unscoped()
		}

		// select
		onlys := req.GetOnly()
		if len(onlys) > 0 {
			db = db.Select(onlys)
		}
		// sort
		if req.OrderBy != "" {
			db = db.Order(fmt.Sprintf("%s %s", req.OrderBy, string(req.OrderDirection)))
		} else {
			db = db.Order(fmt.Sprintf("id %s", string(req.OrderDirection)))
		}
		// search
		if req.Search != "" {
			db = db.Where(fmt.Sprintf("%s LIKE ?", req.GetSearchField()), "%"+req.Search+"%")
		}
		// filters

		// custom query

		//req.Filters.Unmarshal()
		page := req.GetPage()
		if page != -1 {
			db.Limit(req.GetPerPage()).Offset(page - 1*req.GetPerPage())
		}

		models := this.newModels()
		find := db.Find(models)
		if find.Error != nil {
			result <- repository.QueryPaginationResult{Error: find.Error}
			return
		}
		var total int64
		res := find.Count(&total)
		if res.Error != nil {
			result <- repository.QueryPaginationResult{Error: res.Error}
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
		i := reflect.ValueOf(models).Elem().Interface()
		result <- repository.QueryPaginationResult{Result: i, Pagination: pagination}
	}()
	return result
}

func (this *mysqlRep) Create(ctx context.Context, entity interface{}) <-chan repository.InsertResult {
	res := make(chan repository.InsertResult)
	go func() {
		defer close(res)
		create := this.con.Create(entity)
		res <- repository.InsertResult{Id: entity.(model.IModel).GetID(), Result: entity, Error: create.Error}
	}()
	return res
}

func (this *mysqlRep) Save(ctx context.Context, entity interface{}) <-chan repository.QueryResult {
	res := make(chan repository.QueryResult)
	go func() {
		defer close(res)
		saved := this.con.Save(entity)
		res <- repository.QueryResult{Result: entity, Error: saved.Error}
	}()
	return res
}

func (this *mysqlRep) Update(ctx context.Context, id string, update interface{}) <-chan repository.QueryResult {
	res := make(chan repository.QueryResult)
	go func() {
		defer close(res)
		findRes := <-this.FindById(context.Background(), id)
		if findRes.Error != nil {
			res <- repository.QueryResult{Error: findRes.Error}
			return
		}
		saved := this.con.Model(findRes.Result).Updates(update)
		res <- repository.QueryResult{Result: findRes.Result, Error: saved.Error}
	}()
	return res
}

func (this *mysqlRep) Delete(ctx context.Context, id string) <-chan error {
	result := make(chan error)
	go func() {
		defer close(result)
		findRes := <-this.FindById(context.Background(), id)
		if findRes.Error != nil {
			result <- findRes.Error
			return
		}
		deleted := this.con.Delete(findRes.Result)
		result <- deleted.Error
	}()
	return result
}

func (this *mysqlRep) Restore(ctx context.Context, id string) <-chan repository.QueryResult {
	result := make(chan repository.QueryResult)
	go func() {
		defer close(result)
		findRes := <-this.FindById(context.Background(), id)
		if findRes.Error != nil {
			result <- repository.QueryResult{Error: findRes.Error}
			return
		}
		deleted := this.con.Model(findRes.Result).Update("deleted_at", nil)
		result <- repository.QueryResult{Result: findRes.Result, Error: deleted.Error}
	}()
	return result
}
