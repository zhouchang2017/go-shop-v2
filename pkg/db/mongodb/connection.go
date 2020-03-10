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
		Host:       "mongodb-primary",
		Port:       "30000",
		Database:   "go-shop",
		Username:   "root",
		Password:   "12345678",
		AuthSource: "admin",
		ReplicaSet: "rs0",
	})
}

type Config struct {
	Host       string          `json:"host"`
	Port       string          `json:"port"`
	Database   string          `json:"database"`
	Username   string          `json:"username"`
	Password   string          `json:"password"`
	AuthSource string          `json:"auth_source"`
	ReplicaSet string          `json:"replica_set"`
	Ctx        context.Context `json:"-"`
}

func (c Config) URI() string {
	uri := url.URL{}
	uri.Scheme = "mongodb"
	query := url.Values{}
	if c.Host == "" {
		uri.Host = "localhost"
	} else {
		uri.Host = c.Host
		if c.Port != "" {
			uri.Host = fmt.Sprintf("%s:%s", c.Host, c.Port)
		}
	}

	if c.Username != "" && c.Password != "" {
		uri.User = url.UserPassword(c.Username, c.Password)
		authSource := "admin"
		if c.AuthSource != "" {
			authSource = c.AuthSource
		}
		query.Add("authSource", authSource)
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
	log.Printf("Close Mongodb...\n")
}

func CreateIndexes(ctx context.Context, collection *mongo.Collection, models []mongo.IndexModel) (err error) {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err = collection.Indexes().CreateMany(ctx, models, opts)
	if err != nil {
		log.Printf("create indexs error:%s\n", err)
		return err
	}
	return nil
}
