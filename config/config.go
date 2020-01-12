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
}

type config struct {
}

func New() *config {
	return &config{}
}


// rabbitMQ uri
func (c *config) RabbitMQUri() string {
	return "amqp://guest:guest@localhost:5672/"
}

// mongodb config
func (c *config) MongodbConfig() mongodb.Config {
	return mongodb.Config{
		Host:       "localhost",
		Database:   "go-shop",
		Username:   "root",
		Password:   "12345678",
		AuthSource: "go-shop",
	}
}

// mysql config
func (c *config) MysqlConfig() mysql.Config {
	return mysql.Config{
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "go_shop",
		Username: "uhowep",
		Password: "uhowep0770",
	}
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
	return qiniu.Config{
		QiniuDomain:    "q2a066yfd.bkt.clouddn.com",
		QiniuAccessKey: "bZbhwfl0pyHb4EMny9swOtZAhIrJvvzJ7h-NmZaF",
		QiniuSecretKey: "Jtq379mcl9lU0ZOB9rjXmQ_fEZ80fU9G4X3PiEVr",
		Bucket:         "go-shop-v1",
	}
}
