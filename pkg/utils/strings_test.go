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

	spew.Dump(StrPoint("orderCreated"))
	spew.Dump(TodayEnd())
}


