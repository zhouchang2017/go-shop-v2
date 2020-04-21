package http

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/request"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
)

type OrderController struct {
	orderSrv  *services.OrderService
	refundSrv *services.RefundService
	trackRep  *repositories.TrackRep
}

// 退款订单列表
// api GET /refunds
func (ctrl *OrderController) Refunds(ctx *gin.Context) {
	user := ctx2.GetUser(ctx).(*models.User)
	var req request.IndexRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ResponseError(ctx, err)
		return
	}
	req.AppendFilter("user.id", user.GetID())
	req.AppendFilter("refunds", bson.M{"$elemMatch": bson.M{"status": bson.M{"$in": []int{
		models.RefundStatusDone,
		models.RefundStatusAgreed,
		models.RefundStatusApply,
		models.RefundStatusRefunding,
		models.RefundStatusReject}}}})

	orders, pagination, err := ctrl.orderSrv.Pagination(ctx, &req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	// return
	Response(ctx, gin.H{
		"data":       orders,
		"pagination": pagination,
	}, http.StatusOK)

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

	if order, err = ctrl.orderSrv.Cancel(ctx, order, "用户取消"); err != nil {
		ResponseError(ctx, err)
		return
	}
	// 用户取消订单
	rabbitmq.Dispatch(events.NewOrderClosedByUserEvent(order.OrderNo))
	Response(ctx, order, http.StatusOK)
}

// 确认收货
// api PUT /orders/:id/confirm
func (ctrl *OrderController) Confirm(ctx *gin.Context) {
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
	order, err = ctrl.orderSrv.Confirm(ctx, order.OrderNo)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, order, http.StatusOK)
}

// 评价
// api POST /orders/:id/comment
func (ctrl *OrderController) Comment(ctx *gin.Context) {
	var form services.OrderCommentOption
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err2.Err422.F("提交评价失败"))
		return
	}
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
	order, err = ctrl.orderSrv.Comment(ctx, order, user, &form)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, order, http.StatusOK)
}

type refundOption struct {
	Desc string `json:"desc"`
}

// 未发货订单申请退款
// api POST /orders/:id/refunds
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
	// 校验订单状态是否允许用户发起退款
	if err := order.CanApplyRefund(); err != nil {
		ResponseError(ctx, err)
		return
	}

	// 当前只做整单退款
	// 生成退款参数
	opt := &services.MakeRefundOption{
		Desc:  form.Desc,
		Order: order,
	}
	var items []*services.MakeRefundItemOption
	for _, item := range order.OrderItems {
		items = append(items, &services.MakeRefundItemOption{
			ItemId:       item.Item.Id,
			RefundAmount: item.TotalAmount,
			Count:        item.Count,
		})
	}
	opt.Items = items
	refund, err := ctrl.refundSrv.MakeRefund(ctx, opt)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	rabbitmq.Dispatch(events.NewOrderApplyRefundEvent(refund))
	Response(ctx, refund, http.StatusOK)
}

// 获取订单关联退款
// api GET /orders/:id/refunds
func (ctrl *OrderController) OrderRefunds(ctx *gin.Context) {
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
	refunds, err := ctrl.refundSrv.FindRefundByOrderId(ctx, id)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, refunds, http.StatusOK)
}

// 取消退款
// api PUT /refunds/:id/cancel
func (ctrl *OrderController) CancelRefund(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ResponseError(ctx, err2.Err422.F("refund_id缺少"))
		return
	}

	refund, err := ctrl.refundSrv.FindById(ctx, id)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	order, err := ctrl.orderSrv.FindById(ctx, refund.OrderId)
	if err != nil {
		ResponseError(ctx, err2.Err422.F("订单不存在"))
		return
	}
	user := ctx2.GetUser(ctx).(*models.User)
	if order.User.Id != user.GetID() {
		ResponseError(ctx, err2.Err422.F("订单不存在"))
		return
	}
	refund, err = ctrl.refundSrv.Cancel(ctx, refund, user)

	if err != nil {
		ResponseError(ctx, err)
		return
	}

	Response(ctx, refund, http.StatusOK)
}

// 跟踪物流查询
// api GET /tracks/:deliveryId/:wayBillId
func (ctrl *OrderController) Track(ctx *gin.Context) {
	deliveryId := ctx.Param("deliveryId")
	if deliveryId == "" {
		ResponseError(ctx, err2.Err422.F("查询物流失败 [deliveryId]is required"))
		return
	}
	wayBillId := ctx.Param("wayBillId")
	if wayBillId == "" {
		ResponseError(ctx, err2.Err422.F("查询物流失败 [wayBillId]is required"))
		return
	}
	one, err := ctrl.trackRep.FindOne(ctx, bson.M{
		"delivery_id": deliveryId,
		"way_bill_id": wayBillId,
	})
	if err != nil {
		ResponseError(ctx, err)
	}
	Response(ctx, one, http.StatusOK)
}
