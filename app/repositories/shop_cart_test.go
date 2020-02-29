package repositories

import (
	"context"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/db/mongodb"
	"testing"
)

func TestShopCartRep_Count(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	redis.TestConnect()
	defer redis.Close()

	rep := NewShopCartRep(mongodb.GetConFn())
	//rep.cache = redis.GetConFn()

	count := rep.Count(context.Background(), "123")
	t.Log(count)
}
