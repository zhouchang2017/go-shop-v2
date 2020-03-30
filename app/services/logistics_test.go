package services

import (
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/wechat"
	"testing"
)

func TestLogisticsService_GetAllDelivery(t *testing.T) {
	redis.TestConnect()
	defer redis.Close()

	config := wechat.Config{
		AppId:     "wxc55788a7da4b4bfc",
		AppSecret: "4e6f4527564ccee6ef3ea278dc2fc8ef",
	}
	wechat.NewSDK(config)

	service := NewLogisticsService()

	delivery, err := service.GetAllDelivery()
	if err!=nil {
		panic(err)
	}
	spew.Dump(delivery)
}
