package wechat

import (
	"github.com/davecgh/go-spew/spew"
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
