package http

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v2"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/wechat"
	"strings"
)

var wxSrv *weapp.Server
var trackSrv *services.LogisticsService

func registerListeners() {
	// 物流更新监听
	wxSrv.OnExpressPathUpdate(func(result *weapp.ExpressPathUpdateResult) {
		if err := trackSrv.UpdateTrack(context.Background(), result); err != nil {
			// 更新异常
			spew.Dump(err)
		}
	})
}

type WechatController struct {
}

var token = "Gq2l8hEc3wLqMNzn0tABsYzoLZLmRMdj"

func checkSignature(ctx *gin.Context) (string, bool) {
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	echostr := ctx.Query("echostr")
	nonce := ctx.Query("nonce")
	tmpArr := []string{timestamp, nonce, token}
	hash := sha1.New()
	hash.Write([]byte(strings.Join(tmpArr, "")))
	sum := hash.Sum(nil)
	hashString := fmt.Sprintf("%x", sum)
	return echostr, hashString == signature
}

func (this *WechatController) Handle(ctx *gin.Context) {

	if err := wxSrv.Serve(ctx.Writer, ctx.Request); err != nil {
		spew.Dump(err)
	}
}

// 生成小程序码
// api POST /wechat/unlimited/qr-code
type getWechatQrCodeOption struct {
	Scene     string `json:"scene"`
	Page      string `json:"page"`
	Width     int    `json:"width"`
	IsHyaline bool   `json:"is_hyaline"`
}

func (this *WechatController) CreateQrCode(ctx *gin.Context) {
	var form getWechatQrCodeOption
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}
	if form.Width > 1280 {
		form.Width = 1280
	}
	bytes, err := wechat.SDK.UnlimitedQRCode(weapp.UnlimitedQRCode{
		Scene:     form.Scene,
		Page:      form.Page,
		Width:     form.Width,
		AutoColor: true,
		IsHyaline: form.IsHyaline,
	})
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	ctx.Writer.WriteString(string(bytes))
}
