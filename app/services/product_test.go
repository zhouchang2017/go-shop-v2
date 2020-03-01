package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/utils"
	"sync"
	"testing"
)

func TestProductService_List(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeProductService()

	products, pagination, err := service.List(context.Background(), &request.IndexRequest{})

	spew.Dump(products)
	spew.Dump(pagination)

	if err != nil {
		t.Fatal(err)
	}

}

// 随机添加排序值
func TestSort(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeProductService()

	all := <-service.rep.FindAll(context.Background())
	if all.Error != nil {
		t.Fatal(all.Error)
	}
	var wg sync.WaitGroup
	for _, product := range all.Result.([]*models.Product) {
		wg.Add(1)
		go func(p *models.Product) {
			defer wg.Done()
			p.Sort = utils.RandomInt(1000)
			saved := <-service.rep.Save(context.Background(), p)
			if saved.Error != nil {
				t.Fatal(saved.Error)
			}
		}(product)
	}
	wg.Wait()
}

func TestProductService_FindById(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	redis.TestConnect()
	defer redis.Close()

	service := MakeProductService()

	id, err := service.FindById(context.Background(), "5e577e370d3f4744961cfcfd")
	if err!=nil {
		t.Fatal(err)
	}
	spew.Dump(id)
}
