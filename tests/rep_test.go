package tests

import (
	"github.com/davecgh/go-spew/spew"
	models2 "go-shop-v2/app/models"
	"reflect"
	"testing"
)

func TestRefl(t *testing.T) {

	makeModels:= make([]*models2.Product,0)
	spew.Dump(&makeModels)


	modelType := reflect.TypeOf(&models2.Product{})

	slice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0)

	spew.Dump(reflect.New(slice.Type()).Interface())
}
