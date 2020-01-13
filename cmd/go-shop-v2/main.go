package main

import (
	"context"
	"go-shop-v2/app/listeners"
	"go-shop-v2/app/repositories"
	vue2 "go-shop-v2/app/vue"
	"go-shop-v2/config"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/db/mysql"
	"go-shop-v2/pkg/message"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/vue/core"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	config := config.NewConfig()
	// 消息队列
	mq := message.New(config.RabbitMQUri())
	defer mq.Close()
	// 七牛云存储
	newQiniu := qiniu.NewQiniu(config.QiniuConfig())
	// mongodb
	mongoConnect := mongodb.Connect(config.MongodbConfig())
	defer mongodb.Close()
	// mysql
	mysql.Connect(config.MysqlConfig())
	defer mysql.Close()
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
