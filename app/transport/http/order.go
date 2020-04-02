package http

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/request"
	"net/http"
	"strconv"
)

type OrderController struct {
	orderSrv  *services.OrderService
	refundSrv *services.RefundService
}

// 订单列表
func (ctrl *OrderController) Index(ctx *gin.Context) {
	user := ctx2.GetUser(ctx).(*models.User)
	var req request.IndexRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ResponseError(ctx, err)
		return
	}
	if status := ctx.Query("status"); status != "" {
		var value int
		atoi, err := strconv.Atoi(status)
		if err == nil {
			value = atoi
		}
		req.AppendFilter("status", value)
	}

	req.AppendFilter("user.id", user.GetID())

	orders, pagination, err := ctrl.orderSrv.Pagination(ctx, &req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	if len(orders) == 0 {
		// 默认空数组
		orders = []*models.Order{}
	}
	// return
	Response(ctx, gin.H{
		"data":       orders,
		"pagination": pagination,
	}, http.StatusOK)
}

// 订单详情
func (ctrl *OrderController) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	topic, err := ctrl.orderSrv.FindById(ctx, id)
	if err != nil {
		// err
		spew.Dump(err)
	}
	Response(ctx, topic, http.StatusOK)
}

// 创建订单
func (ctrl *OrderController) Store(ctx *gin.Context) {
	// get user information with auth
	userInfo := ctx2.GetUser(ctx).(*models.User)
	// check form data
	var form services.OrderCreateOption
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}

	// create order
	order, err := ctrl.orderSrv.Create(ctx, userInfo, &form)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	// 新订单事件
	rabbitmq.Dispatch(events.NewOrderCreatedEvent(order))
	// 订单延时关闭
	rabbitmq.Dispatch(events.NewOrderTimeOutEvent(order.OrderNo))
	Response(ctx, order, http.StatusOK)
}

// 查询订单状态
// api GET /orders/:id/status
func (ctrl *OrderController) Status(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ResponseError(ctx, err2.Err422.F("缺少订单id参数"))
		return
	}
	status, err := ctrl.orderSrv.GetOrderStatus(ctx, id)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, status, http.StatusOK)
}

// 取消订单
// api PUT /orders/:id/cancel
func (ctrl *OrderController) Cancel(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ResponseError(ctx, err2.Err422.F("订单号异常"))
		return
	}

	order, err := ctrl.orderSrv.FindById(ctx, id)
	if err != nil {
		ResponseError(ctx, err2.Err422.F("订单不存在"))
		return
	}
	user := ctx2.GetUser(ctx).(*models.User)
	if order.User.Id != user.GetID() {
		ResponseError(ctx, err2.Err422.F("订单不存在"))
		return
	}

	if _, err := ctrl.orderSrv.Cancel(ctx, order, "用户取消"); err != nil {
		ResponseError(ctx, err)
		return
	}
	// 用户取消订单
	rabbitmq.Dispatch(events.NewOrderClosedByUserEvent(order.OrderNo))
	Response(ctx, nil, http.StatusNoContent)
}

type refundOption struct {
	Desc string `json:"desc"`
}

// 未发货订单申请退款
// api PUT /orders/:id/refund
func (ctrl *OrderController) ApplyRefund(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ResponseError(ctx, err2.Err422.F("订单号异常"))
		return
	}
	var form refundOption
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}
	order, err := ctrl.orderSrv.FindById(ctx, id)
	if err != nil {
		ResponseError(ctx, err2.Err422.F("订单不存在"))
		return
	}
	user := ctx2.GetUser(ctx).(*models.User)
	if order.User.Id != user.GetID() {
		ResponseError(ctx, err2.Err422.F("订单不存在"))
		return
	}
	refund, order, err := ctrl.orderSrv.ApplyRefund(ctx, order, form.Desc)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	// todo 申请退款成功，推送邮件通知管理员
	rabbitmq.Dispatch(events.NewOrderApplyRefundEvent(order, refund.Id))
	Response(ctx, refund, http.StatusOK)
}

// 取消退款
// api PUT /orders/:id/refund/:refundId/cancel
func (ctrl *OrderController) CancelRefund(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ResponseError(ctx, err2.Err422.F("订单号异常"))
		return
	}

	order, err := ctrl.orderSrv.FindById(ctx, id)
	if err != nil {
		ResponseError(ctx, err2.Err422.F("订单不存在"))
		return
	}
	user := ctx2.GetUser(ctx).(*models.User)
	if order.User.Id != user.GetID() {
		ResponseError(ctx, err2.Err422.F("订单不存在"))
		return
	}
	refundId := ctx.Param("refundId")
	if refundId == "" {
		ResponseError(ctx, err2.Err422.F("退款单号异常"))
		return
	}
	refund, order, err := ctrl.refundSrv.CancelRefund(ctx, &services.RefundOption{
		OrderId:  id,
		OrderNo:  order.OrderNo,
		RefundNo: refundId,
	}, user, true)

	if err != nil {
		ResponseError(ctx, err)
		return
	}

	Response(ctx, refund, http.StatusOK)
}
