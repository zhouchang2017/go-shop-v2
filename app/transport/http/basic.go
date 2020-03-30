package http

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	err2 "go-shop-v2/pkg/err"
)

var guard string

func SetGuard(name string) {
	guard = name
}

func Register(app *gin.Engine) {
	productSrv := services.MakeProductService()
	v1 := app.Group("v1")
	indexController := &IndexController{
		productSrv: productSrv,
		topicSrv:   services.MakeTopicService(),
		articleSrv: services.MakeArticleService(),
	}
	authController := &AuthController{
		userSrv: services.MakeUserService(),
	}

	// 支付
	paymentController := &PaymentController{paymentSrv: services.MakePaymentService(), refundSrv: services.MakeRefundService()}
	// 支付通知回调
	v1.Any("/wechat/payments/paid/notify", paymentController.PayNotify)
	// 退款通知回调
	v1.Any("/wechat/payments/refund/notify", paymentController.RefundNotify)

	// 授权
	v1.POST("/login", authController.Login)

	// 首页列表
	v1.GET("/index", indexController.Index)

	// 文章详情
	v1.GET("/articles/:id", indexController.article)

	// 话题详情
	v1.GET("/topics/:id", indexController.Topic)

	// 话题产品分页
	v1.GET("/topics/:id/products", indexController.TopicProducts)

	productController := &ProductController{
		productSrv:   services.MakeProductService(),
		promotionSrv: services.MakePromotionService(),
	}

	// 产品详情
	v1.GET("/products/:id", productController.Show)

	// 产品促销计划
	v1.GET("/products/:id/promotions", productController.Promotion)

	// 需要授权路由
	v1.Use(auth.AuthMiddleware(guard))
	shopCartController := &ShopCartController{
		srv: services.MakeShopCartService(),
	}
	// 购物车列表页
	v1.GET("/shopping-cart", shopCartController.Index)
	// 加入购物车
	v1.POST("/shopping-cart", shopCartController.Add)
	// 更新购物车
	v1.PUT("/shopping-cart/:id", shopCartController.Update)
	// 更新购物车选定状态
	//v1.PUT("/shopping-cart", shopCartController.UpdateChecked)
	// 删除购物车
	v1.DELETE("/shopping-cart", shopCartController.Delete)
	// 获取购物车选定产品详情
	v1.PUT("/shopping-cart", shopCartController.GetCheckedItemsDetail)

	orderController := &OrderController{
		orderSrv: services.MakeOrderService(),
	}
	// 订单列表
	v1.GET("/orders", orderController.Index)
	// 订单详情
	v1.GET("/orders/:id", orderController.Show)
	// 下单
	v1.POST("/orders", orderController.Store)
	// 查询订单状态
	v1.GET("/orders/:id/status", orderController.Status)
	// 取消订单
	v1.PUT("/orders/:id/cancel", orderController.Cancel)
	// 订单申请退款
	v1.PUT("/orders/:id/refund", orderController.ApplyRefund)
	bookmarkSrv := services.MakeBookmarkService()
	bookmarkController := &BookmarkController{
		bookmarkSrv: bookmarkSrv,
		productSrv:  productSrv,
	}
	// 收藏夹
	// 收藏夹列表页
	v1.GET("/bookmarks", bookmarkController.Index)
	// 当前产品是否被收藏
	v1.GET("/products/:id/bookmarks", bookmarkController.Show)
	// 加入收藏夹
	v1.POST("/products/:id/bookmarks", bookmarkController.Add)
	// 从收藏夹移除
	v1.DELETE("/products/:id/bookmarks", bookmarkController.Delete)

	addressController := &AddressController{addressSrv: services.MakeAddressService()}
	// 用户地址
	// 用户地址列表
	v1.GET("/addresses", addressController.Index)

	// 新增地址
	v1.POST("/addresses", addressController.Add)

	// 更新地址
	v1.PUT("/addresses/:id", addressController.Update)

	// 删除地址
	v1.DELETE("/addresses/:id", addressController.Delete)

	// 统一下单
	v1.POST("/payments", paymentController.UnifiedOrder)

}

//
func Response(ctx *gin.Context, data interface{}, code int) {
	ctx.JSON(code, data)
}

func ResponseXML(ctx *gin.Context, data interface{ ToXmlString() string }, code int) {
	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	ctx.Status(code)
	ctx.Writer.Write([]byte(data.ToXmlString()))
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

	ctx.JSON(errStatus.Code(), errStatus)
	ctx.Abort()
}
