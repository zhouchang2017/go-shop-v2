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

// mongodb config
func (c *configServiceProvider) mongodbConfig() mongodb.Config {
	return mongodb.Config{
		Host:     "localhost",
		Database: "go-shop",
		Username: "root",
		Password: "12345678",
	}
}

// mysql config
func (c *configServiceProvider) mysqlConfig() mysql.Config {
	return mysql.Config{
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "go-shop",
		Username: "root",
		Password: "12345678",
	}
}

// auth config
func (c *configServiceProvider) authGuard(adminRep *repositories.AdminRep) func() auth.StatefulGuard {
	return func() auth.StatefulGuard {
		return auth.NewJwtGuard(
			"admin",
			"admin-secret-key",
			auth.NewRepositoryUserProvider(adminRep),
		)
	}
}

// qiniu config
func (c *configServiceProvider) qiniuConfig() qiniu.Config {
	return qiniu.Config{
		QiniuDomain:    "q2a066yfd.bkt.clouddn.com",
		QiniuAccessKey: "bZbhwfl0pyHb4EMny9swOtZAhIrJvvzJ7h-NmZaF",
		QiniuSecretKey: "Jtq379mcl9lU0ZOB9rjXmQ_fEZ80fU9G4X3PiEVr",
		Bucket:         "go-shop-v1",
	}
}
