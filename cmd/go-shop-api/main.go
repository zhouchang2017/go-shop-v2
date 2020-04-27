package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-shop-v2/app/email"
	"go-shop-v2/app/lbs"
	"go-shop-v2/app/listeners"
	"go-shop-v2/app/services"
	http2 "go-shop-v2/app/transport/http"
	"go-shop-v2/config"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/log"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/wechat"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"
)

var configPathFlag = flag.String("c", ".config", "get the file path for config to parsed")

const PORT = 8081

func main() {

	// parse flag
	flag.Parse()

	// get config path
	if *configPathFlag == "" {
		fmt.Errorf("please use -c to set the config file path or use -h to see more")
		return
	}
	// open file
	file, openErr := os.Open(*configPathFlag)
	if openErr != nil {
		fmt.Errorf("open config file failed caused of %s", openErr.Error())
		return
	}
	// decode json
	decoder := json.NewDecoder(file)

	decodeErr := decoder.Decode(&config.Config)

	file.Close()

	if decodeErr != nil {
		fmt.Errorf("decode config file failed caused of %s", decodeErr.Error())
		return
	}

	configs := config.NewConfig()

	// 邮件服务
	mail := email.New(configs.EmailCfg)

	getwd, _ := os.Getwd()
	join := path.Join(getwd, "runtime", "logs", "go-shop-api.log")
	influxdbConf := configs.InfluxDbCfg
	if influxdbConf != nil {
		influxdbConf.AppName = "api"
	}
	// 日志设置
	log.Setup(log.Option{
		AppName:        "go-shop-api",
		Path:           join,
		MaxAge:         time.Hour * 24 * 30,
		RotationTime:   time.Hour * 24,
		Email:          mail,
		To:             "zhouchangqaz@gmail.com",
		InfluxDBConfig: influxdbConf,
	})

	// 消息队列
	mq := rabbitmq.New(configs.RabbitmqCfg)

	// 七牛云存储
	qiniu.NewQiniu(configs.QiniuConfig())
	// newQiniu := qiniu.NewQiniu(configs.QiniuConfig())
	// mongodb
	mongodb.Connect(configs.MongodbConfig())
	defer mongodb.Close()

	// redis
	redis.Connect(configs.RedisConfig())
	defer redis.Close()

	listeners.ListenerInit()

	// 微信skd
	wechat.NewSDK(configs.WeappConfig)
	wechat.ClearCache() // 清除缓存
	// 微信支付
	wechat.NewPay(configs.WechatPayCfg)

	// 地址解析
	lbs.NewSDK(configs.LbsKey)

	guard := "user"
	// auth service
	authSrv := auth.NewAuth()
	// 注册guard
	authSrv.Register(func() auth.StatefulGuard {
		return auth.NewJwtGuard(
			guard,
			"user-secret-key",
			60*24*7,
			services.MakeUserService(),
		)
	})

	app := gin.New()
	app.Use(log.Logger())
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
	logrus.Infof("[info] start http server listening %d", PORT)

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	ctx2, cancelFunc := context.WithCancel(context.Background())
	mq.RunProducer(ctx2)
	defer mq.Shutdown()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer cancelFunc()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Infof("Server exiting")
}
