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
	productSrv := services.MakeProductService()
	v1 := app.Group("v1")
	indexController := &IndexController{
		productSrv:   productSrv,
		topicSrv:     services.MakeTopicService(),
		articleSrv:   services.MakeArticleService(),
		inventorySrv: services.MakeInventoryService(),
	}
	authController := &AuthController{
		userSrv: services.MakeUserService(),
	}

	// 授权
	v1.POST("/login", authController.Login)

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

	// 需要授权路由
	v1.Use(auth.AuthMiddleware(guard))

	shopCartController := &ShopCartController{
		srv:     services.MakeShopCartService(),
		itemSrv: services.MakeItemService(),
	}
	// 购物车列表页
	v1.GET("/shopping-cart", shopCartController.Index)

	// 加入购物车
	v1.POST("/shopping-cart", shopCartController.Add)

	// 更新购物车
	v1.PUT("/shopping-cart/:id", shopCartController.Update)

	// 更新购物车选定状态
	v1.PUT("/shopping-cart", shopCartController.UpdateChecked)

	// 删除购物车
	v1.DELETE("/shopping-cart", shopCartController.Delete)

	bookmarkSrv := services.MakeBookmarkService()
	bookmarkController := &BookmarkController{
		bookmarkSrv: bookmarkSrv,
		productSrv:  productSrv,
	}
	// 收藏夹
	// 收藏夹列表页
	v1.GET("/bookmarks", bookmarkController.Index)

	// 加入收藏夹
	v1.POST("/bookmarks", bookmarkController.Add)

	// 从收藏夹移除
	v1.DELETE("/bookmarks", bookmarkController.Delete)
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
	ctx.Abort()
}
