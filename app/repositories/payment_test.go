package repositories

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/utils"
	"testing"
)

func TestPaymentRep_GetRangPaymentCount(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	rep := MakePaymentRep()
	start:= utils.StrToTime("2020-03-01 00:00:00")
	end:=utils.StrToTime("2020-03-28 00:00:00")
	count, err := rep.GetRangePaymentCount(context.Background(), start, end)
	if err!=nil {
		panic(err)
	}
	spew.Dump(count)
}

func TestPaymentRep_GetRangPaymentCounts(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	rep := MakePaymentRep()
	start:= utils.StrToTime("2020-03-01 00:00:00")
	end:=utils.StrToTime("2020-03-28 00:00:00")
	count, err := rep.GetRangePaymentCounts(context.Background(), start, end)
	if err!=nil {
		panic(err)
	}
	spew.Dump(count)
}