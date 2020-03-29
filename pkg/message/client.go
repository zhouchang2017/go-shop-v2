package message

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// 事件（Producer）
type Event interface {
	ExchangeName() string
	ExchangeType() string
	RoutingKey() string
	Body() []byte
	Delay() time.Duration
}

// 监听者（Consumer）
type Listener interface {
	Event() Event
	QueueName() string
	OnError(err error)
	Handler(data []byte) error
}

var once sync.Once
var rabbitMQ *RabbitMQ

type RabbitMQ struct {
	wg        sync.WaitGroup
	conn      *amqp.Connection
	listeners []Listener
}

func New(uri string) *RabbitMQ {
	once.Do(func() {
		rabbitMQ = &RabbitMQ{}
		rabbitMQ.connect(uri)
	})
	return rabbitMQ
}

func Dispatch(event Event) {
	rabbitMQ.Dispatch(event)
}

// 注册监听者
func (this *RabbitMQ) Register(listener Listener) {
	this.listeners = append(this.listeners, listener)
}

// 触发事件
func (this *RabbitMQ) Dispatch(event Event) {
	var expiration string
	var routingKey = event.RoutingKey()
	var exchangeName = event.ExchangeName()
	// 是否为延时事件
	if event.Delay().Seconds() > 0 {
		expiration = strconv.Itoa(int(event.Delay().Seconds() * 1000))
		routingKey = this.delayEventRoutingKey(event)
		exchangeName = ""
	}

	channel, err := this.conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel,Error: %s\n", err)
		return
	}

	defer channel.Close()

	err = channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        event.Body(),
			Expiration:  expiration,
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Println("=========Dispatch Event=========")
	log.Printf("Event Name:%s\n", this.eventName(event))
	log.Printf("Delay:%s\n", event.Delay())
	log.Printf("exchange Name:%s\n", exchangeName)
	log.Printf("Routing Key:%s\n", routingKey)
	log.Printf("Payload:%s\n", event.Body())
	log.Printf("\n")
}

func (this *RabbitMQ) eventName(event Event) string {
	t := reflect.TypeOf(event)
	split := strings.Split(t.String(), ".")
	name := split[len(split)-1]
	return name
}

func (this *RabbitMQ) delayEventRoutingKey(event Event) string {
	return fmt.Sprintf("%s_delay", event.RoutingKey())
}

// 初始化
func (this *RabbitMQ) initialize(ctx context.Context) {
	var events []Event

	var inEvents = func(e Event) bool {
		for _, event := range events {
			if event.ExchangeName() == e.ExchangeName() {
				return true
			}
		}
		return false
	}

	for _, listener := range this.listeners {
		if !inEvents(listener.Event()) {
			events = append(events, listener.Event())
		}
	}
	// 注册交换机
	for _, event := range events {
		channel, err := this.conn.Channel()

		if err != nil {
			log.Printf("Failed to open a channel,Error:%s\n", err)
			return
		}

		err = channel.ExchangeDeclare(
			event.ExchangeName(),
			event.ExchangeType(),
			true,
			false,
			false,
			false,
			nil,
		)
		failOnError(err, "Failed to declare an exchange")

		log.Println("=========Queue Declare=========")
		log.Printf("exchange Name:%s\n", event.ExchangeName())
		log.Printf("exchange Type:%s\n", event.ExchangeType())
		log.Printf("\n")

		// 注册延时队列
		if event.Delay().Seconds() > 0 {

			/**
			 * 注意,这里是重点!!!!!
			 * 声明一个延时队列, ß我们的延时消息就是要发送到这里
			 */
			_, errDelay := channel.QueueDeclare(
				this.delayEventRoutingKey(event), // name
				true,                             // durable
				false,                            // delete when unused
				true,                             // exclusive
				false,                            // no-wait
				amqp.Table{
					// 当消息过期时把消息发送到这个 exchange
					"x-dead-letter-exchange": event.ExchangeName(),
				},                                // arguments
			)
			failOnError(errDelay, "Failed to declare a delay_queue")
			log.Println("=========Delay Queue Declare=========")
			log.Printf("Event Name:%s\n", this.eventName(event))
			log.Printf("x-dead-letter-exchange:%s\n", event.ExchangeName())
			log.Printf("Delay:%s\n", event.Delay())
			log.Printf("Routing Key:%s\n", this.delayEventRoutingKey(event))
			log.Printf("\n")

		}

	}
	// 注册队列
	for _, listener := range this.listeners {
		go this.listen(ctx, listener)
	}

}

// 监听
func (this *RabbitMQ) listen(ctx context.Context, listener Listener) {
	channel, err := this.conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel,Error:%s\n", err)
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				channel.Close()
				log.Printf("close channel...")
				return
			}
		}
	}()

	queue, err := channel.QueueDeclare(
		listener.QueueName(),
		true,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = channel.QueueBind(
		queue.Name,
		listener.Event().RoutingKey(),
		listener.Event().ExchangeName(),
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	// 获取消费通道
	channel.Qos(1, 0, true) // 确保rabbitmq会一个一个发消息
	msgs, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	go func() {

		for msg := range msgs {

			select {
			case <-ctx.Done():
				log.Printf("listener ctx down!")
				channel.Close()
				log.Printf("close channel...")
				return

			default:
				if err := listener.Handler(msg.Body); err != nil {
					listener.OnError(err)
					msg.Ack(false)
					return
				}
				msg.Ack(true)
			}

		}

	}()

}

func (this *RabbitMQ) connect(uri string) {
	connection, err := amqp.Dial(uri)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect rabbitMQ,err:%s\n", err))
	}
	this.conn = connection
	log.Printf("Connection rabbitMQ [%s]\n", uri)
}

func (this *RabbitMQ) Run(ctx context.Context) {
	go this.initialize(ctx)
}

func (this *RabbitMQ) Close() {
	log.Printf("Close rabbitMQ\n", )
	this.conn.Close()
}
