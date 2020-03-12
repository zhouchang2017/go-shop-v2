package repositories

import (
	"context"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
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

func (this *ItemRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "code", Value: bsonx.Int64(-1)}},
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
	}
}

func NewItemRep(rep repository.IRepository) *ItemRep {
	repository := &ItemRep{rep}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Printf("model [%s] create indexs error:%s\n", repository.TableName(), err)
		panic(err)
	}
	return repository
}
