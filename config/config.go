package config

import (
	"context"
	"github.com/Abramovic/logrus_influxdb"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go-shop-v2/app/email"
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/db/mysql"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/rabbitmq"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/wechat"
	"os"
	"sync"
)

func init() {
	fields.DefaultMapLocation = &fields.MapValue{
		Lng: 112.969625,
		Lat: 28.199554,
	}
}

var Config *config
var once sync.Once
var Path string

type config struct {
	v            *viper.Viper
	LbsKey       string                  `json:"lbs_key" mapstructure:"lbs_key"`
	WeappConfig  wechat.Config           `json:"weapp_config" mapstructure:"weapp_config"`
	WechatPayCfg wechat.PayConfig        `json:"wechatpay_config" mapstructure:"wechatpay_config"`
	MongoTestCfg mongodb.Config          `json:"mongo_config_test" mapstructure:"mongo_config_test"`
	MongoCfg     mongodb.Config          `json:"mongo_config" mapstructure:"mongo_config"`
	MysqlCfg     mysql.Config            `json:"mysql_config" mapstructure:"mysql_config"`
	RedisCfg     redis.Config            `json:"redis_config" mapstructure:"redis_config"`
	QiniuCfg     qiniu.Config            `json:"qiniu_config" mapstructure:"qiniu_config"`
	EmailCfg     email.Config            `json:"email_config" mapstructure:"email_config"`
	RabbitmqCfg  rabbitmq.Config         `json:"rabbitmq_config" mapstructure:"rabbitmq_config"`
	InfluxDbCfg  *logrus_influxdb.Config `json:"influxdb_config" mapstructure:"influxdb_config"`
}

func New() *config {
	once.Do(func() {
		Config = &config{}
		Config.init()
		Config.load()
	})
	return Config
}

func (c *config) init() {
	c.v = viper.New()
	c.v.SetConfigName(".config")
	c.v.SetConfigType("json")
	if Path == "" {
		Path, _ = os.Getwd()
	}
	c.v.AddConfigPath(Path)
	c.v.AddConfigPath("/etc/go-shop")
}

func (c *config) load() error {
	if err := c.v.ReadInConfig(); err != nil {
		return err
	}
	if err := c.v.Unmarshal(&c); err != nil {
		return err
	}
	return nil
}

func (c *config) Watch() error {
	ctx, cancel := context.WithCancel(context.Background())
	c.v.WatchConfig()
	//监听回调函数
	watch := func(e fsnotify.Event) {
		c.load()
		cancel()
	}
	c.v.OnConfigChange(watch)
	<-ctx.Done()
	return nil
}

func NewConfig() *config {
	return Config
}

// wechat payment config
func (c *config) WechatPayConfig() wechat.PayConfig {
	return c.WechatPayCfg
}

// mongodb config
func (c *config) MongodbConfig() mongodb.Config {
	return c.MongoCfg
}

// mysql config
func (c *config) MysqlConfig() mysql.Config {
	return c.MysqlCfg
}

// redis config
func (c *config) RedisConfig() redis.Config {
	return c.RedisCfg
}

//// auth config
//func (c *config) authGuard(adminRep *repositories.AdminRep) func() auth.StatefulGuard {
//	return func() auth.StatefulGuard {
//		return auth.NewJwtGuard(
//			"admin",
//			"admin-secret-key",
//			auth.NewRepositoryUserProvider(adminRep),
//		)
//	}
//}

// qiniu config
func (c *config) QiniuConfig() qiniu.Config {
	return c.QiniuCfg
}
