package redis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"sync"
)

var once sync.Once
var con *Connection

var GetConFn func() *Connection = func() *Connection {
	return con
}

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
	//Prefix   string `json:"prefix"`
}

func (config Config) Options() *redis.Options {
	var host = "localhost"
	if config.Host != "" {
		host = config.Host
	}
	var port = "6379"
	if config.Port != "" {
		port = config.Port
	}
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: config.Password,
		DB:       config.Database,
	}
}

type Connection struct {
	*redis.Client
}

func TestConnect() {
	Connect(Config{
		Host: "localhost",
		Port: "63790",
	})
}

func Connect(c Config) *Connection {
	once.Do(func() {
		con = &Connection{redis.NewClient(c.Options())}
	})
	return con
}

func Close() {
	con.Close()
}
