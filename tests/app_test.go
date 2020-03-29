package tests

import (
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/rabbitmq/app"
	"testing"
)

func Test_time(t *testing.T)  {
	config := app.Config{
		Username: "root",
		Password: "12345678",
	}

	spew.Dump(config.URI())
}
