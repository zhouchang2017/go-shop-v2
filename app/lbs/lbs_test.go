package lbs

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestNewSDK(t *testing.T) {
	sdk := NewSDK("K3DBZ-27U6P-M34DP-V2RF3-6BHI7-AWF3O")

	res, err := sdk.ParseAddress("湖南省株洲市炎陵县政府")
	if err!=nil {
		t.Fatal(err)
	}
	spew.Dump(res)
}
