package http

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/wechat"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	ctx2 "go-shop-v2/pkg/ctx"
	"net/http"
)

type PaymentController struct {
	paymentSrv *services.PaymentService
}

// 统一下单
// api /payments
func (p *PaymentController) UnifiedOrder(ctx *gin.Context) {
	var form services.PaymentOption
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}
	form.SpbillCreateIp = ctx.ClientIP()
	user := ctx2.GetUser(ctx).(*models.User)

	wechatMiniPayConfig, err := p.paymentSrv.Payment(ctx, user, &form)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, wechatMiniPayConfig, http.StatusOK)
}

// 回调
func (p *PaymentController) PayNotify(ctx *gin.Context) {
	err := p.paymentSrv.PayNotify(ctx, ctx.Request)
	spew.Dump("支付回调异常:")
	spew.Dump(err)
	rsp := new(wechat.NotifyResponse) // 回复微信的数据
	if err != nil {
		rsp.ReturnCode = gopay.FAIL
		rsp.ReturnMsg = gopay.FAIL
		ResponseXML(ctx, rsp, http.StatusOK)
		return
	}
	rsp.ReturnCode = gopay.SUCCESS
	rsp.ReturnMsg = gopay.OK
	ResponseXML(ctx, rsp, http.StatusOK)
	return
}
