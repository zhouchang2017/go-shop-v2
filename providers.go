package go_shop_v2

import (
	"go-shop-v2/app/listeners"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/app/vue/resources"
	"go-shop-v2/config"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/db"
	"go-shop-v2/pkg/event"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/support"
	"go-shop-v2/pkg/transport/http"
	"go-shop-v2/pkg/vue"
)

var Providers []func() support.ServiceProvider = []func() support.ServiceProvider{
	event.NewEventServiceProvider,
	config.NewConfigServiceProvider,
	listeners.NewListenerServiceProvider,
	qiniu.NewQiniuServiceProvider,
	http.NewHttpServiceProvider,
	db.NewMongodbServiceProvider,
	db.NewMysqlServiceProvider,
	auth.NewAuthServiceProvider,
	vue.NewVueServiceProvider,
	models.NewModelServiceProvider,
	repositories.NewRepositoryServiceProvider,
	services.NewServiceServiceProvider,
	resources.NewVueResourceServiceProvider,
}

var TestProvider []func() support.ServiceProvider = []func() support.ServiceProvider{
	event.NewEventServiceProvider,
	config.NewConfigServiceProvider,
	listeners.NewListenerServiceProvider,
	qiniu.NewQiniuServiceProvider,
	//http.NewHttpServiceProvider,
	db.NewMongodbServiceProvider,
	db.NewMysqlServiceProvider,
	auth.NewAuthServiceProvider,
	//vue.NewVueServiceProvider,
	models.NewModelServiceProvider,
	repositories.NewRepositoryServiceProvider,
	services.NewServiceServiceProvider,
	resources.NewVueResourceServiceProvider,
}