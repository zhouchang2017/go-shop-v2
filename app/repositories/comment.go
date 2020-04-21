package repositories

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type CommentRep struct {
	repository.IRepository
}

func (this *CommentRep) CreateMany(ctx context.Context, comments []*models.Comment) error {
	entities := []interface{}{}
	for _, comment := range comments {
		if comment.ID.IsZero() {
			comment.ID = primitive.NewObjectID()
			comment.CreatedAt = time.Now()
			comment.UpdatedAt = time.Now()
		}
		entities = append(entities, comment)
	}
	if len(entities) == 0 {
		return errors.New("评价内容不能为空")
	}
	_, err := this.Collection().InsertMany(ctx, entities)
	if err != nil {
		return err
	}
	return nil
}

func (this *CommentRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "order_no", Value: bsonx.Int64(-1)},
				{Key: "product_id", Value: bsonx.Int64(-1)},
			}, // order_no 唯一键
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
	}
}

func NewCommentRep(rep repository.IRepository) *CommentRep {
	repository := &CommentRep{rep}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Panicf("model [%s] create indexs error:%s\n", repository.TableName(), err)
	}
	return repository
}
