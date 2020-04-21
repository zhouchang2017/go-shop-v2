package wechat

import (
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/cache/redis"
	"testing"
	"time"
)

func TestSdk_GetAllDelivery(t *testing.T) {
	redis.TestConnect()
	defer redis.Close()

	config := Config{
		AppId:     "",
		AppSecret: "",
	}
	newSDK := NewSDK(config)
	delivery, err := newSDK.GetAllDelivery()
	if err != nil {
		panic(err)
	}
	spew.Dump(delivery)
}

func TestSdk_GetAllAccount(t *testing.T) {
	redis.TestConnect()
	defer redis.Close()

	config := Config{
		AppId:     "",
		AppSecret: "",
	}
	newSDK := NewSDK(config)
	accounts, err := newSDK.GetAllAccount()
	if err != nil {
		panic(err)
	}
	spew.Dump(accounts)
}

func TestSdk_TestUpdateOrder(t *testing.T) {
	redis.TestConnect()
	defer redis.Close()

	config := Config{
		AppId:     "",
		AppSecret: "",
	}
	newSDK := NewSDK(config)

	order, err := newSDK.TestUpdateOrder(&TestExpressUpdateOrderOption{
		OrderId:    "1250955380666667008",
		WaybillID:  "1250955380666667008_1587090440_waybill_id",
		ActionTime: time.Now(),
		ActionType: 300003,
		ActionMsg:  "您的快件已签收，如有疑问请电联小哥【赖辅勇，电话：19859307707】。疫情期间顺丰每日对网点消毒、小哥每日测温、配戴口罩，感谢您使用顺丰，期待再次为您服务。",
	})
	if err!=nil {
		panic(err)
	}
	spew.Dump(order)
}
