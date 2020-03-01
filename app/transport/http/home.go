package http

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
)

type HomeController struct {
	bookmarkSrv *services.BookmarkService
}

func (this *HomeController) Index(ctx *gin.Context) {
	// 用户home页面数据
	// 积分、收藏夹、优惠券
}

// 用户地址列表
func (this *HomeController) Addresses(ctx *gin.Context) {

}
