package http

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	ctx2 "go-shop-v2/pkg/ctx"
	"net/http"
)

type PaymentController struct {
	paymentSrv *services.PaymentService
}

// 下单
// api /checkouts
func (p *PaymentController) CheckOut(ctx *gin.Context) {
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
