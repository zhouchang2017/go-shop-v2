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
	srv *services.ShopCartService
}

type shopCartForm struct {
	ItemId string `json:"item_id" form:"item_id"`
	Qty    uint64  `json:"qty"`
}

// 个人购物车列表
func (this *ShopCartController) Index(ctx *gin.Context) {
	user := ctx2.GetUser(ctx).(*models.User)
	req := &request.IndexRequest{}

	shopCartItems, pagination, err := this.srv.Index(ctx, user.GetID(), req.GetPage(), req.GetPerPage())
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	Response(ctx, gin.H{
		"data":       shopCartItems,
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
	user := ctx2.GetUser(ctx).(*models.User)

	if form.Qty == 0 {
		form.Qty = 1
	}

	item, err := this.srv.Add(ctx, user.GetID(), form.ItemId, form.Qty)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	Response(ctx, item, 200)
}

type updateShopCartForm struct {
	ItemId *string `json:"item_id" form:"item_id"`
	Qty    uint64   `json:"qty"`
}

// 更新购物车
// 若传入ItemId，则是在购物车页面点击option值，进行更新
// 否则只是对购物车商品数量进行增减
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
	user := ctx2.GetUser(ctx).(*models.User)

	// itemId != nil 在购物车点击更新
	if form.ItemId != nil {
		status, shopCartItem, err := this.srv.Update(ctx, user.GetID(), id, *form.ItemId, form.Qty)
		if err != nil {
			ResponseError(ctx, err)
			return
		}
		Response(ctx, gin.H{
			"status": status,
			"data":   shopCartItem,
		}, http.StatusOK)
		return
	}

	shopCartItem, err := this.srv.UpdateQty(ctx, user.GetID(), id, form.Qty)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	Response(ctx, gin.H{
		"status": 3,
		"data":   shopCartItem,
	}, 200)
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
	user := ctx2.GetUser(ctx).(*models.User)
	if err := this.srv.Toggle(ctx, user.GetID(), form.Checked, form.Ids...); err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, form.Ids, 200)
}

// 获取选定items详情
func (this *ShopCartController) GetCheckedItemsDetail(ctx *gin.Context) {
	var form updateShopCartCheckedForm
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}
	user := ctx2.GetUser(ctx).(*models.User)
	detail, err := this.srv.GetShopCartItemsDetail(ctx, user.GetID(), form.Ids...)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, detail, 200)
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
	user := ctx2.GetUser(ctx).(*models.User)
	if err := this.srv.Delete(ctx, user.GetID(), form.Ids...); err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, form.Ids, 200)
}
