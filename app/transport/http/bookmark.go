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

type BookmarkController struct {
	bookmarkSrv *services.BookmarkService
	productSrv  *services.ProductService
}

// 个人收藏夹列表
func (this *BookmarkController) Index(ctx *gin.Context) {
	user := ctx2.GetUser(ctx).(*models.User)

	var req request.IndexRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	bookmarks, pagination, err := this.bookmarkSrv.Index(ctx, user.GetID(), req.Page, req.PerPage)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ids, err := this.productSrv.FindByIds(ctx, bookmarks)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	Response(ctx, gin.H{
		"data":       ids,
		"pagination": pagination,
	}, http.StatusOK)
}

type bookmarkForm struct {
	ProductId string `json:"product_id" form:"product_id"`
}

// 添加到收藏夹
func (this *BookmarkController) Add(ctx *gin.Context) {
	productId := ctx.Param("id")

	if productId == "" {
		ResponseError(ctx, err2.New(http.StatusUnprocessableEntity, "product_id参数缺少"))
		return
	}

	user := ctx2.GetUser(ctx).(*models.User)

	err := this.bookmarkSrv.Add(ctx, user.GetID(), productId)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, nil, http.StatusNoContent)
}

type deleteBookmarkForm struct {
	Ids []string `json:"ids"`
}

// 从收藏夹移除
func (this *BookmarkController) Delete(ctx *gin.Context) {
	productId := ctx.Param("id")

	if productId == "" {
		ResponseError(ctx, err2.New(http.StatusUnprocessableEntity, "product_id参数缺少"))
		return
	}
	user := ctx2.GetUser(ctx).(*models.User)

	if err := this.bookmarkSrv.Remove(ctx, user.GetID(), productId); err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, nil, http.StatusNoContent)
}

// 当前产品收藏夹状态
func (this *BookmarkController) Show(ctx *gin.Context) {
	user := ctx2.GetUser(ctx).(*models.User)
	productId := ctx.Param("id")
	if productId == "" {
		ResponseError(ctx, err2.New(422, "缺少product_id"))
		return
	}
	bookmark := this.bookmarkSrv.FindByProductId(ctx, user.GetID(), productId)
	isCollect:= true
	if bookmark == nil {
		isCollect = false
	}
	Response(ctx, isCollect, http.StatusOK)
}
