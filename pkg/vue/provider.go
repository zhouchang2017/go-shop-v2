package vue

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
	"log"
	"net/http"
	"time"
)

const port = ":8083"

type vueServiceProvider struct {
	app    *gin.Engine
	server *http.Server
}

func NewVueServiceProvider() support.ServiceProvider {
	return &vueServiceProvider{}
}

func (v *vueServiceProvider) Register(container support.Container) {
	container.Provide(NewVue)
}

func (v *vueServiceProvider) Boot() fx.Option {
	return fx.Invoke(v.start)
}

func (v *vueServiceProvider) start(lifecycle fx.Lifecycle, vue *Vue) {
	v.app = gin.New()
	v.app.Use(gin.Logger())
	v.app.Use(gin.Recovery())
	v.app.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": port})
	})

	group := v.app.Group("app")


	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			vue.HttpHandler(group)
			v.server = &http.Server{
				Addr:           port,
				Handler:        v.app,
				ReadTimeout:    time.Second * 30,
				WriteTimeout:   time.Second * 30,
				MaxHeaderBytes: 1 << 20,
			}
			log.Printf("[info] start http server listening %s", port)
			go func() {
				// 连接服务器
				if err := v.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(i context.Context) error {
			log.Println("Shutdown Server ...")
			ctx, cancel := context.WithTimeout(i, 5*time.Second)
			defer cancel()
			return v.server.Shutdown(ctx)
		},
	})
}
