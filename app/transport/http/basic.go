package http

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
)

func Register(app *gin.Engine) {

	v1 := app.Group("v1")
	productController := &ProductController{
		productSrv: services.MakeProductService(),
		topicSrv:   services.MakeTopicService(),
		articleSrv: services.MakeArticleService(),
	}

	v1.GET("/index", productController.Index)

}

//
func Response(ctx *gin.Context, data interface{}, code int) {
	ctx.JSON(code, data)
}

// 错误响应
func ResponseError(ctx *gin.Context, err error) {

}
