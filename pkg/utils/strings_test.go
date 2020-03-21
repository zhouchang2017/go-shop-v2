package utils

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/inflection"
	"testing"
	"time"
)

func TestStructToName(t *testing.T) {

	//name := StrToPlural("Brand")
	singular := StrToSnake(inflection.Singular("BrandOptions"))

	t.Logf(singular)
}

func TestStrToPlural(t *testing.T) {

	spew.Dump(time.Now().Format("20060102150405"))
}

