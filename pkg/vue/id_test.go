package vue

import (
	"github.com/davecgh/go-spew/spew"
	uuid "github.com/satori/go.uuid"
	"go-shop-v2/app/models"
	"reflect"
	"testing"
)

func TestNewIDField(t *testing.T) {
	field := NewIDField()

	data := struct {
		ID string `json:"id"`
	}{
		ID:uuid.NewV4().String(),
	}

	field.Resolve(nil,data)


	spew.Dump(field)
}

func TestNewTextField(t *testing.T) {
	//field := NewTextField("省份", "Address.Province")
	//spew.Dump(field)

	var addr = models.Shop{
		Name:"123123",
		Address:&models.ShopAddress{
			Addr:"中国",
			Areas:"深圳",
		},
	}
	r := reflect.ValueOf(addr)
	f := reflect.Indirect(r).FieldByName("Address.Addr")


	spew.Dump(f)
}