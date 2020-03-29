package mp_subscribe

import (
	"context"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/wechat"
	"testing"
)

func TestNewOrderPaidNotify(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	wechat.NewSDK(wechat.Config{
		AppId:     "wxc55788a7da4b4bfc",
		AppSecret: "4e6f4527564ccee6ef3ea278dc2fc8ef",
	})

	service := services.MakeOrderService()
	order, err := service.FindById(context.Background(), "5e7da762f1309bd23bb714cc")
	if err != nil {
		panic(err)
	}
	err = wechat.SDK.SendSubscribeMessage(NewOrderPaidNotify(order))
	if err != nil {
		t.Fatal(err)
	}
}
