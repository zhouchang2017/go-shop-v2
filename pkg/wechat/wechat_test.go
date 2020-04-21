package wechat

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/medivhzhan/weapp/v2"
	"go-shop-v2/pkg/cache/redis"
	"testing"
	"time"
)

func TestSdk_GetDailySummary(t *testing.T) {
	redis.TestConnect()
	defer redis.Close()

	config := Config{
		AppId:     "",
		AppSecret: "",
	}
	newSDK := NewSDK(config)

	summary, err := newSDK.GetDailyVisitTrend(time.Now())
	if err != nil {
		panic(err)
	}
	spew.Dump(summary)
}

func TestSdk_UnlimitedQRCode(t *testing.T) {
	redis.TestConnect()
	defer redis.Close()

	config := Config{
		AppId:     "",
		AppSecret: "",
	}
	newSDK := NewSDK(config)
	code, err := newSDK.UnlimitedQRCode(weapp.UnlimitedQRCode{
		Scene:     "id=5e858e0019cadb2de0130492",
		Page:      "pages/product",
		Width:     0,
		AutoColor: true,
		IsHyaline: true,
	})
	if err!=nil {
		panic(err)
	}

	spew.Dump(string(code))
}