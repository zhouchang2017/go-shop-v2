package db

import (
	"context"
	"github.com/jinzhu/gorm"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/db/mysql"
	"go-shop-v2/pkg/support"
	"go.uber.org/fx"
	"log"
)

func NewMysqlServiceProvider() support.ServiceProvider {
	return &mysqlServiceProvider{}
}

type mysqlServiceProvider struct {
}

func (m *mysqlServiceProvider) Register(container support.Container) {
	container.Provide(mysql.Connect)
}

func (m *mysqlServiceProvider) Boot() fx.Option {
	return fx.Invoke(m.close)
}

func (m *mysqlServiceProvider) close(lifecycle fx.Lifecycle, con *gorm.DB) {
	lifecycle.Append(fx.Hook{
		OnStop: func(context context.Context) error {
			log.Print("mysql closing")
			return con.Close()
		},
	})
}

func NewMongodbServiceProvider() support.ServiceProvider {
	return &mongodbServiceProvider{}
}

type mongodbServiceProvider struct {
}

func (m *mongodbServiceProvider) Register(container support.Container) {
	container.Provide(mongodb.Connect)
}

func (m *mongodbServiceProvider) close(lifecycle fx.Lifecycle, con *mongodb.Connection) {
	lifecycle.Append(fx.Hook{
		OnStop: func(context context.Context) error {
			log.Print("mongodb closing")
			return con.Client().Disconnect(context)
		},
	})
}

func (m *mongodbServiceProvider) Boot() fx.Option {
	return fx.Invoke(m.close)
}
