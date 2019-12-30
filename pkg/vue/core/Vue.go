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
	app        *gin.Engine
	server     *http.Server
	port       int64
	prefix     string
	resources  []contracts.Resource
	guard      string
	httpHandle *httpHandle
	warps      []*warp
	customHttpHandle []func(router gin.IRouter)
}

func New(port int64) *Vue {
	engine := gin.New()
	return &Vue{
		app:    engine,
		port:   port,
		prefix: "app",
	}
}

func (this *Vue) init() {
	this.app.Use(gin.Logger())
	this.app.Use(gin.Recovery())
	group := this.app.Group(this.prefix)
	this.setWarps()
	this.httpHandle = newHttpHandle(this, group)
	this.httpHandle.exec()
}

func (this *Vue) setWarps() {
	for _, resource := range this.resources {
		this.warps = append(this.warps, newWarp(resource))
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
