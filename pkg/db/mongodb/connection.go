package mongodb

import (
	"context"
	"fmt"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/url"
	"reflect"
	"sync"
	"time"
)

var once sync.Once
var con *Connection

var GetConFn func() *Connection = func() *Connection {
	return con
}

func TestConnect() {
	Connect(Config{
		Host:     "localhost",
		Database: "go-shop",
		Username: "root",
		Password: "12345678",
	})
}

type Config struct {
	Host       string
	Database   string
	Username   string
	Password   string
	AuthSource string
	ReplicaSet string
	Ctx        context.Context
}

func (c Config) URI() string {
	uri := url.URL{}
	uri.Scheme = "mongodb"
	query := url.Values{}
	if c.Host == "" {
		uri.Host = "localhost"
	} else {
		uri.Host = c.Host
	}

	if c.Username != "" && c.Password != "" {
		uri.User = url.UserPassword(c.Username, c.Password)
		query.Add("authSource", c.AuthSource)
	}
	if c.Database != "" {
		uri.Path = fmt.Sprintf("/%s", c.Database)
	}
	if c.ReplicaSet != "" {
		query.Add("replicaSet", c.ReplicaSet)
	}
	uri.RawQuery = query.Encode()
	res := uri.String()
	defer log.Printf("mongodb URI = %s", res)
	return res
}

func Connect(c Config) *Connection {
	once.Do(func() {
		opts := options.Client()
		opts.ApplyURI(c.URI())
		var ctx context.Context = c.Ctx

		if c.Ctx == nil {
			ctx = context.Background()
		}

		if c.ReplicaSet != "" {
			opts.SetReadPreference(readpref.Primary())
			opts.SetServerSelectionTimeout(time.Duration(2 * time.Second))
		}

		opts.SetMaxPoolSize(5)

		ctxTimeout, _ := context.WithTimeout(ctx, 10*time.Second)
		client, err := mongo.Connect(ctxTimeout, opts)
		if err != nil {
			panic(err)
		}

		timeout, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelFunc()
		err = client.Ping(timeout, readpref.Primary())
		if err != nil {
			panic(err)
		}
		con = &Connection{client: client, database: c.Database}
	})
	return con
}

type Connection struct {
	client   *mongo.Client
	database string
}

func (c *Connection) Client() *mongo.Client {
	return c.client
}

func (c *Connection) Database() *mongo.Database {
	return c.client.Database(c.database)
}

func (c *Connection) Collection(model interface{}) *mongo.Collection {
	var name string
	if reflect.TypeOf(model).Kind() == reflect.String {
		name = model.(string)
	} else {
		name = utils.StructNameToSnakeAndPlural(model)
	}
	return c.Database().Collection(name)
}

func Close() {
	if err := con.client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
