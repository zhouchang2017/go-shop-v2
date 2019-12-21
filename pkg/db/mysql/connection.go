package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
)

var migrates []interface{}
var con *gorm.DB
var once sync.Once

var GetConFn func() *gorm.DB = func() *gorm.DB {
	return con
}

func Register(entity interface{}) {
	migrates = append(migrates, entity)
}

func TestConnect() {
	config := Config{
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "go-shop",
		Username: "root",
		Password: "123456",
	}

	_ = Connect(config)
}

type Config struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func (c Config) URI() string {
	uri := url.URL{}
	query := url.Values{}
	if c.Host != "" {
		uri.Host = c.Host
	}

	if c.Port != "" {
		uri.Host = fmt.Sprintf("(%s)", net.JoinHostPort(c.Host, c.Port))
	}
	if c.Username != "" && c.Password != "" {
		uri.User = url.UserPassword(c.Username, c.Password)
	}
	if c.Database != "" {
		uri.Path = c.Database
	}
	query.Add("charset", "utf8mb4")
	query.Add("parseTime", "True")
	query.Add("loc", "Local")
	uri.RawQuery = query.Encode()

	s := strings.TrimPrefix(uri.String(), "//")
	defer log.Printf("mysql URI = %s", s)
	return s
}

func Connect(conf Config) *gorm.DB {
	once.Do(func() {
		db, err := gorm.Open("mysql", conf.URI())
		if err != nil {
			panic(err)
		}
		con = db

		// run migrates
		con.AutoMigrate(migrates...)
	})
	return con
}

func Close() {
	if err := con.Close(); err != nil {
		panic(err)
	}
}
