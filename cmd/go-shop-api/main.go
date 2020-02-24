package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/repositories"
	http2 "go-shop-v2/app/transport/http"
	"go-shop-v2/config"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/message"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/wechat"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var configPathFlag = flag.String("c", ".env", "get the file path for config to parsed")

const PORT = 8081

func main() {
	// parse flag
	flag.Parse()

	// get config path
	if *configPathFlag == "" {
		fmt.Println("please use -c to set the config file path or use -h to see more")
		return
	}
	// open file
	file, openErr := os.Open(*configPathFlag)
	if openErr != nil {
		fmt.Println("open config file failed caused of %s", openErr.Error())
		return
	}
	// decode json
	decoder := json.NewDecoder(file)

	decodeErr := decoder.Decode(&config.Config)

	file.Close()

	if decodeErr != nil {
		fmt.Printf("decode config file failed caused of %s", decodeErr.Error())
		return
	}

	configs := config.NewConfig()
	// 消息队列
	mq := message.New(configs.RabbitMQUri())
	defer mq.Close()
	// 七牛云存储
	qiniu.NewQiniu(configs.QiniuConfig())
	// newQiniu := qiniu.NewQiniu(configs.QiniuConfig())
	// mongodb
	mongoConnect := mongodb.Connect(configs.MongodbConfig())
	defer mongodb.Close()

	// 微信skd
	wechat.NewSDK(configs.WeappConfig)


	guard:="user"
	// auth service
	authSrv := auth.NewAuth()
	// 注册guard
	authSrv.Register(func() auth.StatefulGuard {
		return auth.NewJwtGuard(
			guard,
			"user-secret-key",
			repositories.NewUserRep(mongoConnect),
		)
	})

	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	http2.SetGuard(guard)

	http2.Register(app)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", PORT),
		Handler:        app,
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("[info] start http server listening %d", PORT)

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
