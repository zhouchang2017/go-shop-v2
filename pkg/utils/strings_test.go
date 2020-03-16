package utils

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/inflection"
	"math"
	"testing"
)

func TestStructToName(t *testing.T) {

	//name := StrToPlural("Brand")
	singular := StrToSnake(inflection.Singular("BrandOptions"))

	t.Logf(singular)
}

func TestStrToPlural(t *testing.T) {

	type item struct {
		qty           int64
		price         int64
		salePrice     int64
		unitSalePrice int64
	}

	data := []*item{
		{
			qty:   1,
			price: 47500,
		},
		{
			qty:   2,
			price: 53300,
		},
	}

	var salePrices int64 = 5300
	var totalPrice int64
	for _, i := range data {
		totalPrice += i.qty * i.price
	}

	avg := float64(salePrices) / float64(totalPrice)

	var used int64
	for index, i := range data {
		if index+1 == len(data) {
			i.salePrice = salePrices - used
			i.unitSalePrice = int64(math.Ceil(float64(i.salePrice) / float64(i.qty)))
		} else {
			i.unitSalePrice = int64(math.Ceil(avg * float64(i.price)))
			i.salePrice = i.unitSalePrice * i.qty
			used += i.salePrice
		}
	}
	spew.Dump(data)
}
