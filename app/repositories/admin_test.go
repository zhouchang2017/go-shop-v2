package repositories

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestAdminRep_FindByNotifies(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	rep := MakeAdminRep()
	notifies, err := rep.FindByNotifies(context.Background(), []string{"OrderPaid"}, options.Find().SetProjection(bson.M{"email": 1, "nickname": 1}))
	if err!=nil {
		panic(err)
	}
	spew.Dump(notifies)
}
