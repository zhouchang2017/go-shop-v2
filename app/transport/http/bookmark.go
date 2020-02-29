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
	user := ctx2.GetUser(ctx)
	currentUser := user.(*models.User)

	var req request.IndexRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	req.AppendFilter("user_id", currentUser.GetID())
	bookmarks, pagination, err := this.bookmarkSrv.Pagination(ctx, &req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, gin.H{
		"data":       bookmarks,
		"pagination": pagination,
	}, http.StatusOK)
}

type bookmarkForm struct {
	ProductId string `json:"product_id" form:"product_id"`
}

// 添加到收藏夹
func (this *BookmarkController) Add(ctx *gin.Context) {
	var form bookmarkForm
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}
	if form.ProductId == "" {
		ResponseError(ctx, err2.New(http.StatusUnprocessableEntity, "product_id参数缺少"))
		return
	}

	product, err := this.productSrv.FindById(ctx, form.ProductId)
	if err != nil {
		// 产品不存在
		ResponseError(ctx, err)
		return
	}
	user := ctx2.GetUser(ctx)
	currentUser := user.(*models.User)
	add, err := this.bookmarkSrv.Add(ctx, product, currentUser.GetID())
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, add, http.StatusOK)
}

type deleteBookmarkForm struct {
	Ids []string `json:"ids"`
}

// 从购物车移除
func (this *BookmarkController) Delete(ctx *gin.Context) {
	var form deleteBookmarkForm
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}

	if err := this.bookmarkSrv.Delete(ctx, form.Ids...); err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, form.Ids, 200)
}
