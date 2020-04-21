package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-shop-v2/app/email"
	"go-shop-v2/app/lbs"
	"go-shop-v2/app/listeners"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	vue2 "go-shop-v2/app/vue"
	"go-shop-v2/config"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/log"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/wechat"
	"os"
	"os/signal"
	"path"
	"time"
)

var configPathFlag = flag.String("c", ".config", "get the file path for config to parsed")

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
		fmt.Printf("decode config file failed caused of %s", decodeErr.Error())
		return
	}

	configs := config.NewConfig()
	// 邮件服务
	mail := email.New(configs.EmailCfg)

	getwd, _ := os.Getwd()
	join := path.Join(getwd, "runtime", "logs", "go-shop-backend.log")
	// 日志设置
	log.Setup(log.Option{
		AppName:      "go-shop-backend",
		Path:         join,
		MaxAge:       time.Hour * 24 * 30,
		RotationTime: time.Hour * 24,
		Email:        mail,
		To:           "zhouchangqaz@gmail.com",
	})

	// 消息队列
	mq := rabbitmq.New(configs.RabbitmqCfg)

	// 七牛云存储
	newQiniu := qiniu.NewQiniu(configs.QiniuConfig())

	fields.DefaultFileUploadAction = newQiniu.FileUploadAction()

	// mongodb
	mongoConnect := mongodb.Connect(configs.MongodbConfig())
	defer mongodb.Close()
	// mysql
	//mysql.Connect(configs.MysqlConfig())
	//defer mysql.Close()
	// redis
	connect := redis.Connect(configs.RedisConfig())
	defer redis.Close()

	// 刷新缓存
	connect.FlushDB()

	listeners.ListenerInit()

	// 库存初始化
	// inventoryRep := repositories.NewInventoryRep(repositories.NewBasicMongoRepositoryByDefault(&models.Inventory{}, mongodb.GetConFn()))
	// inventoryRep.InitCache()

	// 微信skd
	wechat.NewSDK(configs.WeappConfig)
	wechat.ClearCache() // 清除缓存
	// 微信支付
	wechat.NewPay(configs.WechatPayCfg)

	// 地址解析
	lbs.NewSDK(configs.LbsKey)

	adminGuard := "admin"
	// auth service
	authSrv := auth.NewAuth()
	// 注册guard
	authSrv.Register(func() auth.StatefulGuard {
		return auth.NewJwtGuard(
			adminGuard,
			"admin-secret-key",
			60*24,
			auth.NewRepositoryUserProvider(
				repositories.NewAdminRep(repositories.NewBasicMongoRepositoryByDefault(&models.Admin{}, mongoConnect)),
			),
		)
	})

	// 实例化vue后台组件
	vue := core.New()
	vue.SetLoggerMiddleware(log.Logger())
	// 设置授权守卫
	vue.SetGuard(adminGuard)
	// 注册七牛api
	vue.RegisterCustomHttpHandler(newQiniu.HttpHandle)
	// 注册全局前端配置
	vue.WithConfig("qiniu_upload_action", newQiniu.FileUploadAction())
	vue.WithConfig("qiniu_cdn_domain", newQiniu.Domain())
	vue.WithConfig("logistics", models.LogisticsInfos)
	vue.WithConfig("order_status", models.OrderStatus)
	vue.WithConfig("refund_status", models.RefundStatus)

	notifications := listeners.GetAppNotifications()
	vue.WithConfig("notifications", notifications)
	// vue相关启动项
	vue2.Boot(vue)
	// 启动vue后台组件框架
	vue.Run(8083)

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
	if err := vue.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Infof("Server exiting")
}
