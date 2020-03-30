package wechat

import (
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/cache/redis"
	"testing"
)

func TestSdk_GetAllDelivery(t *testing.T) {
	redis.TestConnect()
	defer redis.Close()

	config := Config{
		AppId:     "wxc55788a7da4b4bfc",
		AppSecret: "4e6f4527564ccee6ef3ea278dc2fc8ef",
	}
	newSDK := NewSDK(config)
	delivery, err := newSDK.GetAllDelivery()
	if err!=nil {
		panic(err)
	}
	spew.Dump(delivery)
}

func TestSdk_GetAllAccount(t *testing.T) {
	redis.TestConnect()
	defer redis.Close()

	config := Config{
		AppId:     "wxc55788a7da4b4bfc",
		AppSecret: "4e6f4527564ccee6ef3ea278dc2fc8ef",
	}
	newSDK := NewSDK(config)
	accounts, err := newSDK.GetAllAccount()
	if err!=nil {
		panic(err)
	}
	spew.Dump(accounts)
}
