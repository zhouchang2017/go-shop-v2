package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestProductRep_Create(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	itemRep := NewItemRep(NewBasicMongoRepositoryByDefault(&models.Item{}, mongodb.GetConFn()))

	productRep := NewProductRep(NewBasicMongoRepositoryByDefault(&models.Product{}, mongodb.GetConFn()), itemRep)

	var session mongo.Session
	var err error
	var ctx context.Context

	ctx = context.Background()
	
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		t.Fatal(err)
	}
	if err = session.StartTransaction(); err != nil {
		t.Fatal(err)
	}

	mongo.WithSession(ctx,session, func(sessionContext mongo.SessionContext) error {
		// 创建两个商品

		product1 := &models.Product{
			Name: utils.RandomOrderNo("测试商品"),
			Code: utils.RandomString(16),
		}

		created := <-productRep.Create(sessionContext, product1)
		if created.Error != nil {
			return created.Error
		}
		product2 := &models.Product{
			Name: utils.RandomOrderNo("测试商品"),
			Code: utils.RandomString(16),
		}
		_ = <-productRep.Create(sessionContext, product2)

		//if created2.Error != nil {
		//	sessionContext.AbortTransaction(sessionContext)
		//	return created2.Error
		//}
		//sessionContext.CommitTransaction(sessionContext)
		session.CommitTransaction(sessionContext)
		return err2.Err401
	})
	session.EndSession(ctx)
}
