package repositories

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
	"testing"
)

func TestBookmarkRep_Add(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()



	rep := NewBookmarkRep(NewBasicMongoRepositoryByDefault(&models.Bookmark{}, mongodb.GetConFn()))

	err := rep.Add(context.Background(), "123", "5e577e370d3f4744961cfcfg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBookmarkRep_Remove(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	rep := NewBookmarkRep(NewBasicMongoRepositoryByDefault(&models.Bookmark{}, mongodb.GetConFn()))

	err := rep.Remove(context.Background(), "123", "5e577e370d3f4744961cfcfd","5e577e370d3f4744961cfcfe")
	if err != nil {
		t.Fatal(err)
	}
}
