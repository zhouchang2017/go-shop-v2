package repositories

import (
	"context"
	"go-shop-v2/pkg/repository"
)

type ItemRep struct {
	repository.IRepository
}

func (this *ItemRep) FindByProductId(ctx context.Context, id string) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	go func() {
		defer close(output)

		many := <-this.FindMany(ctx, map[string]interface{}{"product.id": id})

		output <- repository.QueryResult{
			Result: many.Result,
			Error:  many.Error,
		}
	}()
	return output
}

func NewItemRep(rep repository.IRepository) *ItemRep {
	return &ItemRep{rep}
}
