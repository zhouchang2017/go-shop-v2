package repositories

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type UserRep struct {
	repository.IRepository
}

func (this *UserRep) RetrieveById(identifier interface{}) (auth.Authenticatable, error) {
	result := <-this.FindById(context.Background(), identifier.(string))
	if result.Error != nil {
		return nil, result.Error
	}
	return result.Result.(*models.User), nil
}

func (this *UserRep) RetrieveByCredentials(credentials map[string]string) (auth.Authenticatable, error) {
	openId, ok := credentials["open_id"]
	if !ok && credentials == nil {
		return nil, errors.New("credentials is empty!")
	}
	return this.FindByOpenId(context.Background(), openId)
}

func (this *UserRep) ValidateCredentials(user auth.Authenticatable, credentials map[string]string) bool {
	return true
}

func (this *UserRep) FindByOpenId(ctx context.Context, openId string) (user *models.User, err error) {
	result := this.Collection().FindOne(ctx, bson.M{
		"wechat_mini_id": openId,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	user = &models.User{}
	err = result.Decode(user)
	if err != nil {
		return nil, err
	}
	return
}

func (this *UserRep) index() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "wechat_mini_id", Value: bsonx.Int64(-1)}},
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
	}
}

func NewUserRep(rep repository.IRepository) *UserRep {
	repository := &UserRep{rep}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := repository.Collection().Indexes().CreateMany(context.Background(), repository.index(), opts)
	if err != nil {
		log.Panicf("model [%s] create indexs error:%s\n", repository.TableName(), err)
	}
	return repository
}
