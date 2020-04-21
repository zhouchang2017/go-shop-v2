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
	"go-shop-v2/config"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/log"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/wechat"
	"os"
	"os/signal"
	"path"
	"time"
)

// 消息队列处理进程
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
		fmt.Errorf("decode config file failed caused of %s", decodeErr.Error())
		return
	}

	configs := config.NewConfig()

	// 邮件服务
	mail := email.New(configs.EmailCfg)

	getwd, _ := os.Getwd()
	join := path.Join(getwd, "runtime", "logs", "go-shop-worker.log")
	// 日志设置
	log.Setup(log.Option{
		AppName:      "go-shop-worker",
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
	mongodb.Connect(configs.MongodbConfig())
	defer mongodb.Close()

	// redis
	redis.Connect(configs.RedisConfig())
	defer redis.Close()

	// 微信skd
	wechat.NewSDK(configs.WeappConfig)
	wechat.ClearCache() // 清除缓存

	// 微信支付
	wechat.NewPay(configs.WechatPayCfg)

	// 地址解析
	lbs.NewSDK(configs.LbsKey)
	listeners.ListenerInit()

	listeners.Boot(mq)

	ctx2, cancelFunc := context.WithCancel(context.Background())
	mq.RunConsumer(ctx2)
	defer mq.Shutdown()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Info("Shutdown Server ...")
	defer cancelFunc()

}
