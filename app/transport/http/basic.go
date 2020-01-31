package http

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
)

func Register(app *gin.Engine) {

	v1 := app.Group("v1")
	indexController := &IndexController{
		productSrv: services.MakeProductService(),
		topicSrv:   services.MakeTopicService(),
		articleSrv: services.MakeArticleService(),
	}

	// 首页列表
	v1.GET("/index", indexController.Index)

	// 文章详情
	v1.GET("/articles/:id", indexController.article)

	// 话题详情
	v1.GET("/topics/:id", indexController.Topic)
}

//
func Response(ctx *gin.Context, data interface{}, code int) {
	ctx.JSON(code, data)
}

// 错误响应
func ResponseError(ctx *gin.Context, err error) {

}
