package http

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	err2 "go-shop-v2/pkg/err"
	"net/http"
)

var guard string

func SetGuard(name string) {
	guard = name
}

func Register(app *gin.Engine) {

	v1 := app.Group("v1")
	indexController := &IndexController{
		productSrv:   services.MakeProductService(),
		topicSrv:     services.MakeTopicService(),
		articleSrv:   services.MakeArticleService(),
		inventorySrv: services.MakeInventoryService(),
	}
	authController := &AuthController{
		userSrv: services.MakeUserService(),
	}

	// 授权
	v1.POST("/login", authController.Login)
	// 注册
	v1.POST("/register", authController.Register)

	v1.Use(auth.AuthMiddleware(guard))
	// 首页列表
	v1.GET("/index", indexController.Index)

	// 文章详情
	v1.GET("/articles/:id", indexController.article)

	// 话题详情
	v1.GET("/topics/:id", indexController.Topic)

	// 产品详情
	v1.GET("/products/:id", indexController.Product)

	// 淘宝详情接口
	v1.GET("/taobao/:id", indexController.TaobaoDetail)
}

//
func Response(ctx *gin.Context, data interface{}, code int) {
	ctx.JSON(code, data)
}

// 错误响应
func ResponseError(ctx *gin.Context, err error) {
	var errStatus *err2.Status
	switch err.(type) {
	case *err2.Status:
		errStatus = err.(*err2.Status)
	default:
		errStatus = err2.New(500, err.Error())
	}

	ctx.JSON(http.StatusOK, errStatus)
}
