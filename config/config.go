package config

import (
	"go-shop-v2/pkg/cache/redis"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/db/mysql"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/wechat"
)

func init() {
	fields.DefaultMapLocation = &fields.MapValue{
		Lng: 112.969625,
		Lat: 28.199554,
	}
}

var Config *config

type config struct {
	LbsKey       string         `json:"lbs_key"`
	WeappConfig  wechat.Config  `json:"weapp_config"`
	MongoTestCfg mongodb.Config `json:"mongo_config_test"`
	MongoCfg     mongodb.Config `json:"mongo_config"`
	MysqlCfg     mysql.Config   `json:"mysql_config"`
	RedisCfg     redis.Config   `json:"redis_config"`
	QiniuCfg     qiniu.Config   `json:"qiniu_config"`
}

func NewConfig() *config {
	return Config
}

// rabbitMQ uri
func (c *config) RabbitMQUri() string {
	return "amqp://guest:guest@localhost:5672/"
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
