package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	register(NewItemRep)
}

type ItemRep struct {
	*mongoRep
}

func (this *ItemRep) FindByProductId(ctx context.Context, id string) <-chan repository.QueryResult {
	output := make(chan repository.QueryResult)
	go func() {
		defer close(output)
		productId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			output <- repository.QueryResult{Error: err}
			return
		}
		many := <-this.mongoRep.FindMany(ctx, map[string]interface{}{"product._id": productId})

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
