package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
	"log"
	"net/http"
	"time"
)

var addr string = ":8081"

type httpServiceProvider struct {
	server *http.Server
}

func (h *httpServiceProvider) Register(container support.Container) {
	container.Provide(h.NewGin)
}

func NewHttpServiceProvider() support.ServiceProvider {
	return &httpServiceProvider{}
}

func (h *httpServiceProvider) NewGin() *gin.Engine {
	return gin.New()
}

func (h *httpServiceProvider) start(lifecycle fx.Lifecycle, app *gin.Engine) {
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": addr})
	})
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			h.server = &http.Server{
				Addr:           addr,
				Handler:        app,
				ReadTimeout:    time.Second * 30,
				WriteTimeout:   time.Second * 30,
				MaxHeaderBytes: 1 << 20,
			}
			log.Printf("[info] start http server listening %s", addr)
			go func() {
				// 连接服务器
				if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(i context.Context) error {

			log.Println("Shutdown Server ...")
			ctx, cancel := context.WithTimeout(i, 5*time.Second)
			defer cancel()
			return h.server.Shutdown(ctx)
		},
	})
}

func (h *httpServiceProvider) Boot() fx.Option {
	return fx.Options(fx.Invoke(h.start))
}
