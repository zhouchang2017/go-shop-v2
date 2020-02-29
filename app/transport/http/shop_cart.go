package http

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"net/http"
)

type ShopCartController struct {
	srv     *services.ShopCartService
	itemSrv *services.ItemService
}

type shopCartForm struct {
	ItemId string `json:"item_id" form:"item_id"`
	Qty    int64  `json:"qty"`
}

// 个人购物车列表
func (this *ShopCartController) Index(ctx *gin.Context) {
	user := ctx2.GetUser(ctx)
	currentUser := user.(*models.User)

	req := &request.IndexRequest{}

	req.AppendFilter("user_id", currentUser.GetID())
	carts, pagination, err := this.srv.Pagination(ctx, req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, gin.H{
		"data":       carts,
		"pagination": pagination,
	}, http.StatusOK)
}

// 加入商品到购物车
func (this *ShopCartController) Add(ctx *gin.Context) {
	form := &shopCartForm{}
	if err := ctx.ShouldBind(form); err != nil {
		ResponseError(ctx, err)
		return
	}
	user := ctx2.GetUser(ctx)
	currentUser := user.(*models.User)
	item, err := this.itemSrv.FindById(ctx, form.ItemId)
	if err != nil {
		// 产品不存在
		ResponseError(ctx, err)
		return
	}
	if form.Qty == 0 {
		form.Qty = 1
	}
	cart, err := this.srv.Add(ctx, currentUser.GetID(), item, form.Qty, true)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, cart, 200)
}

type updateShopCartForm struct {
	ItemId  *string `json:"item_id" form:"item_id"`
	Qty     int64   `json:"qty"`
	Checked bool    `json:"checked"`
}

// 更新购物车
func (this *ShopCartController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ResponseError(ctx, err2.New(http.StatusUnprocessableEntity, "缺少id参数"))
		return
	}
	var form updateShopCartForm
	err := ctx.ShouldBind(&form)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	user := ctx2.GetUser(ctx)
	currentUser := user.(*models.User)

	var item *models.Item

	if form.ItemId != nil {
		item, err = this.itemSrv.FindById(ctx, *form.ItemId)
		if err != nil {
			// 产品不存在
			ResponseError(ctx, err)
			return
		}
	}
	forceContext := ctx2.WithForce(ctx, true)
	updated, err := this.srv.Update(forceContext, currentUser.GetID(), id, item, form.Qty, form.Checked)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	Response(ctx, updated, 200)
}

type updateShopCartCheckedForm struct {
	Ids     []string `json:"ids"`
	Checked bool     `json:"checked"`
}

// 更新选定
func (this *ShopCartController) UpdateChecked(ctx *gin.Context) {
	var form updateShopCartCheckedForm
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}
	if err := this.srv.CheckedOrCancelAll(ctx, form.Checked, form.Ids...); err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, form.Ids, 200)
}

type deleteShopCartForm struct {
	Ids []string `json:"ids"`
}

// 删除
func (this *ShopCartController) Delete(ctx *gin.Context) {
	var form deleteShopCartForm
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}
	forceContext := ctx2.WithForce(ctx, true)
	if err := this.srv.Delete(forceContext, form.Ids...); err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, form.Ids, 200)
}
