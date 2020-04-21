package utils

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/inflection"
	"testing"
)

func TestStructToName(t *testing.T) {

	//name := StrToPlural("Brand")
	singular := StrToSnake(inflection.Singular("BrandOptions"))

	t.Logf(singular)
}

func TestStrToPlural(t *testing.T) {
	n := Int64AvgN(500, 3)
	spew.Dump(n)
	randomString := RandomString(32)
	spew.Dump(randomString)
}


