package utils

import (
	"github.com/jinzhu/inflection"
	"testing"
)

func TestStructToName(t *testing.T) {

	//name := StrToPlural("Brand")
	singular := StrToSnake(inflection.Singular("BrandOptions"))

	t.Logf(singular)
}
