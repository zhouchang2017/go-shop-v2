package http

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"net/http"
	"strconv"
)

type OrderController struct {
	orderSrv *services.OrderService
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

	if err := ctrl.orderSrv.Cancel(ctx, order); err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, nil, http.StatusNoContent)
}
