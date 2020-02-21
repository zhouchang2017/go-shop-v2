package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/repository"
)

type ItemRep struct {
	*mongoRep
}

func (this *ItemRep) FindByProductId(ctx context.Context, id string) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	go func() {
		defer close(output)

		many := <-this.mongoRep.FindMany(ctx, map[string]interface{}{"product.id": id})

		output <- repository.QueryResult{
			Result: many.Result,
			Error:  many.Error,
		}
	}()
	return output
}

func NewItemRep(con *mongodb.Connection) *ItemRep {
	return &ItemRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Item{}, con),
	}
}
