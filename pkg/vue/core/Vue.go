package core

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/vue/contracts"
	"log"
	"net/http"
	"time"
)

type Vue struct {
	app              *gin.Engine
	server           *http.Server
	port             int64
	prefix           string
	resources        []contracts.Resource
	pages            []contracts.Page
	guard            string
	httpHandle       *httpHandle
	warps            map[string]*warp
	customHttpHandle []func(router gin.IRouter)
}

func New(port int64) *Vue {
	engine := gin.New()
	return &Vue{
		app:    engine,
		port:   port,
		prefix: "app",
		warps:  map[string]*warp{},
	}
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func (this *Vue) init() {
	this.app.Use(gin.Logger())
	this.app.Use(gin.Recovery())
	this.app.Use(Cors())
	group := this.app.Group(this.prefix)
	this.setWarps()
	this.httpHandle = newHttpHandle(this, group)
	this.httpHandle.exec()
}

func (this *Vue) setWarps() {
	for _, resource := range this.resources {
		this.warps[ResourceUriKey(resource)] = newWarp(resource)
	}
}

func (this *Vue) Run() error {
	this.init()
	this.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", this.port),
		Handler:        this.app,
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("[info] start http server listening %d", this.port)

	go func() {
		// 连接服务器
		if err := this.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return nil
}

func (this *Vue) Shutdown(ctx context.Context) error {
	log.Println("Shutdown Server ...")
	i, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return this.server.Shutdown(i)
}

func (this *Vue) SetPrefix(prefix string) *Vue {
	this.prefix = prefix
	return this
}

// 设置授权守卫
func (this *Vue) SetGuard(guard string) *Vue {
	this.guard = guard
	return this
}

func (this *Vue) RegisterCustomHttpHandler(handler func(router gin.IRouter)) {
	this.customHttpHandle = append(this.customHttpHandle, handler)
}

// 注册资源
func (this *Vue) RegisterResource(resource contracts.Resource) {
	this.resources = append(this.resources, resource)
}

// 注册自定义页面
func (this *Vue) RegisterPage(page contracts.Page) {
	this.pages = append(this.pages, page)
}
