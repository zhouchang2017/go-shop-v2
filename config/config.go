package config

import (
	"encoding/json"
	"fmt"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/db/mysql"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/vue/fields"
	"os"
)

func init() {
	fields.DefaultMapLocation = &fields.MapValue{
		Lng: 112.969625,
		Lat: 28.199554,
	}
}

type config struct {
	MongoCfg mongodb.Config `json:"mongo_config"`
	MysqlCfg mysql.Config   `json:"mysql_config"`
	QiniuCfg qiniu.Config   `json:"qiniu_config"`
}

func NewConfig() *config {
	// todo:fix bug that rewrite this func
	envPath, err := utils.GetFilePath(2, ".env")
	if err != nil {
		panic(fmt.Sprintf("get config path failed caused of %s", err.Error()))
	}
	// open file
	file, openErr := os.Open(envPath)
	if openErr != nil {
		panic(fmt.Sprintf("open config file failed caused of %s", openErr.Error()))
	}
	defer file.Close()
	// decode json
	decoder := json.NewDecoder(file)
	var config config
	decodeErr := decoder.Decode(&config)
	if decodeErr != nil {
		panic(fmt.Sprintf("decode config file failed caused of %s", decodeErr.Error()))
	}
	// return
	return &config
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
