package config

import (
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/db/mysql"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/vue/fields"
)

func init() {
	fields.DefaultMapLocation = &fields.MapValue{
		Lng: 112.969625,
		Lat: 28.199554,
	}
	// 文件默认上传地址
	fields.DefaultFileUploadAction = "https://upload-z2.qiniup.com"
}

var Config *config

type config struct {
	MongoCfg mongodb.Config `json:"mongo_config"`
	MysqlCfg mysql.Config   `json:"mysql_config"`
	QiniuCfg qiniu.Config   `json:"qiniu_config"`
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

// auth config
func (c *config) authGuard(adminRep *repositories.AdminRep) func() auth.StatefulGuard {
	return func() auth.StatefulGuard {
		return auth.NewJwtGuard(
			"admin",
			"admin-secret-key",
			auth.NewRepositoryUserProvider(adminRep),
		)
	}
}

// qiniu config
func (c *config) QiniuConfig() qiniu.Config {
	return c.QiniuCfg
}
