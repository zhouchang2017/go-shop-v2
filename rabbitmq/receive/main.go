package main

import (
	log "github.com/sirupsen/logrus"
	"go-shop-v2/rabbitmq/app"
)

var uri = "amqp://root:12345678@localhost:5672/"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}
func main() {

	app.Register(&app.OnOrderCreatedNotifyAdmin{})
	app.Register(&app.OnOrderTimeoutNotifyAdmin{})
	//app.Register(&app.OnOrderTimoutNotifyUser{})

	rmq, err := app.NewReceive(uri, "go-shop.direct")
	if err != nil {
		log.Fatalf("run: failed to init rabbitmq: %v", err)
	}

	forever := make(chan bool)
	defer rmq.Shutdown()

	<-forever

}
