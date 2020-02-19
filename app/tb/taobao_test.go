package tb

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestTaobaoSdkService_Detail(t *testing.T) {
	service := &TaobaoSdkService{}

	data, err := service.Detail("600740693156")
	if err!=nil {
		t.Fatal(err)
	}

	spew.Dump(data)
}
