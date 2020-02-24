package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go-shop-v2/app/listeners"
	"go-shop-v2/app/repositories"
	vue2 "go-shop-v2/app/vue"
	"go-shop-v2/config"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/message"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/wechat"
	"log"
	"os"
	"os/signal"
	"time"
)

var configPathFlag = flag.String("c", ".env", "get the file path for config to parsed")

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
	newQiniu := qiniu.NewQiniu(configs.QiniuConfig())

	fields.DefaultFileUploadAction = newQiniu.FileUploadAction()

	// mongodb
	mongoConnect := mongodb.Connect(configs.MongodbConfig())
	defer mongodb.Close()
	// mysql
	//mysql.Connect(configs.MysqlConfig())
	//defer mysql.Close()


	// 微信skd
	wechat.NewSDK(configs.WeappConfig)


	// auth service
	authSrv := auth.NewAuth()
	// 注册guard
	authSrv.Register(func() auth.StatefulGuard {
		return auth.NewJwtGuard(
			"admin",
			"admin-secret-key",
			auth.NewRepositoryUserProvider(repositories.NewAdminRep(mongoConnect)),
		)
	})

	// 注册事件监听者
	listeners.Boot(mq)
	// 实例化vue后台组件
	vue := core.New(8083)
	// 设置授权守卫
	vue.SetGuard("admin")
	// 注册七牛api
	vue.RegisterCustomHttpHandler(newQiniu.HttpHandle)
	// 注册全局前端配置
	vue.WithConfig("qiniu_upload_action", newQiniu.FileUploadAction())
	vue.WithConfig("qiniu_cdn_domain", newQiniu.Domain())
	// vue相关启动项
	vue2.Boot(vue)
	// 启动vue后台组件框架
	vue.Run()

	ctx2, cancelFunc := context.WithCancel(context.Background())
	mq.Run(ctx2)

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer cancelFunc()
	if err := vue.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
