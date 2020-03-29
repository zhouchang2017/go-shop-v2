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



	rmq, err := app.NewSender(uri, "go-shop.direct")
	if err != nil {
		log.Fatalf("run: failed to init rabbitmq: %v", err)
	}
	//forever := make(chan bool)
	defer rmq.Shutdown()
	//tick := time.Tick(time.Second * 10)
	//for {
	//	<-tick
	//	err = rmq.Publish("user.event.publish", []byte(utils.RandomOrderNo("非延时任务")), 6000)
	//	if err != nil {
	//		log.Fatalf("run: failed to publish into rabbitmq: %v", err)
	//	}
	//}

	err = rmq.Publish(&app.OrderTimeOutEvent{})
	if err != nil {
		log.Fatalf("run: failed to publish into rabbitmq: %v", err)
	}

	//<-forever

}
