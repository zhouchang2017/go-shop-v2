package wechat

import (
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/utils"
	"testing"
)

func TestPay_UnifiedOrder(t *testing.T) {
	newPay := NewPay(PayConfig{

	})

	order, err := newPay.UnifiedOrder(&PayUnifiedOrderOption{
		Body:           "测试支付",
		OutTradeNo:     utils.RandomOrderNo(""),
		TotalFee:       1,
		SpbillCreateIp: "127.0.0.1",
		OpenId:         "oE1ny5N5B1gbhPlnlgcqI0OS_r-A",
	})
	if err!=nil {
		t.Fatal(err)
	}
	spew.Dump(order)
}
