package rabbitmq

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go-shop-v2/pkg/utils"
	"net/url"
	"sync"
	"time"
)

type Event interface {
	Delay() time.Duration
	Body() []byte
}

type EventRouterKey interface {
	RouterKey() string
}

type Listener interface {
	Make() Listener
	Event() Event
	OnError(payload []byte, err error)
	Handler(data []byte) error
}

var instance *RabbitMQ
var once sync.Once

type RabbitMQ struct {
	config       Config
	exchange     string
	conn         *amqp.Connection
	channel      *amqp.Channel
	closeChannel chan *amqp.Error
	quitChannel  chan bool
	listeners    []Listener
	retryTTL     int
}

func (r RabbitMQ) Listeners() []Listener {
	return r.listeners
}

func (r RabbitMQ) masterExchangeName() string {
	return r.exchange
}

func (r RabbitMQ) delayedExchangeName() string {
	return fmt.Sprintf("%s.delayed", r.exchange)
}

func (r RabbitMQ) retryExchangeName() string {
	return fmt.Sprintf("%s.retry", r.exchange)
}

func durationToInt(d time.Duration) int {
	return int(d.Seconds() * 1000)
}

func Dispatch(event Event) error {
	return instance.Publish(event)
}

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	VHost    string `json:"vhost"`
	Exchange string `json:"exchange"`
}

func (c Config) URI() string {
	uri := url.URL{}
	uri.Scheme = "amqp"
	query := url.Values{}
	host := c.Host
	if c.Host == "" {
		host = "localhost"
	}
	port := c.Port
	if c.Port == "" {
		port = "5672"
	}
	uri.Host = fmt.Sprintf("%s:%s", host, port)

	if c.Username != "" && c.Password != "" {
		uri.User = url.UserPassword(c.Username, c.Password)
	}
	if c.VHost != "" {
		uri.Path = fmt.Sprintf("/%s", c.VHost)
	}

	uri.RawQuery = query.Encode()
	res := uri.String()
	defer log.Printf("rabbitmq URI = %s", res)
	return res
}

func New(c Config) *RabbitMQ {
	once.Do(func() {
		duration := time.Second * 5
		instance = &RabbitMQ{
			exchange: c.Exchange,
			config:   c,
			retryTTL: durationToInt(duration),
		}
	})
	return instance
}

// master 主Exchange，发布消息时发布到该Exchange
// master.delayed 延时Exchange，延时消息发布到该Exchange
// master.retry 重试Exchange，消息处理失败时（3次以内），将消息重新投递给该Exchange

func (rmq *RabbitMQ) RunProducer(ctx context.Context) (*RabbitMQ, error) {
	err := rmq.initProducer()
	if err != nil {
		return nil, err
	}
	rmq.quitChannel = make(chan bool)
	go rmq.handleDisconnect(rmq.initProducer)
	return rmq, err
}

func (rmq *RabbitMQ) Register(listeners ...Listener) *RabbitMQ {
	rmq.listeners = append(rmq.listeners, listeners...)
	return rmq
}

func (rmq *RabbitMQ) RunConsumer(ctx context.Context) (*RabbitMQ, error) {
	err := rmq.initConsumer()
	if err != nil {
		return nil, err
	}

	rmq.quitChannel = make(chan bool)

	go rmq.handleDisconnect(rmq.initConsumer)

	return rmq, err
}

func (rmq *RabbitMQ) initProducer() error {
	var err error

	rmq.conn, err = amqp.Dial(rmq.config.URI())
	if err != nil {
		return err
	}

	rmq.channel, err = rmq.conn.Channel()
	if err != nil {
		return err
	}

	log.Info("connection to rabbitMQ established")

	rmq.closeChannel = make(chan *amqp.Error)
	rmq.conn.NotifyClose(rmq.closeChannel)

	// declare exchange if not exist
	// master 主Exchange，发布消息时发布到该Exchange
	err = rmq.channel.ExchangeDeclare(rmq.masterExchangeName(), "direct", true, false, false, false, nil)
	if err != nil {
		return errors.Wrapf(err, "declaring exchange %q", rmq.masterExchangeName())
	}

	// master.delayed 延时Exchange，延时消息发布到该Exchange
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	if err := rmq.channel.ExchangeDeclare(rmq.delayedExchangeName(), "x-delayed-message", true, false, false, false, args); err != nil {
		return errors.Wrapf(err, "declaring exchange %q", rmq.delayedExchangeName())
	}

	// master.retry 重试Exchange，消息处理失败时（3次以内），将消息重新投递给该Exchange
	if err := rmq.channel.ExchangeDeclare(rmq.retryExchangeName(), "direct", true, false, false, false, nil); err != nil {
		return errors.Wrapf(err, "declaring  exchange %q", rmq.retryExchangeName())
	}
	return nil
}

func (rmq *RabbitMQ) initConsumer() error {
	var err error

	rmq.conn, err = amqp.Dial(rmq.config.URI())
	if err != nil {
		return err
	}

	rmq.channel, err = rmq.conn.Channel()
	if err != nil {
		return err
	}

	log.Info("connection to rabbitMQ established")

	rmq.closeChannel = make(chan *amqp.Error)
	rmq.conn.NotifyClose(rmq.closeChannel)

	// declare exchange if not exist
	// master 主Exchange，发布消息时发布到该Exchange
	err = rmq.channel.ExchangeDeclare(rmq.masterExchangeName(), "direct", true, false, false, false, nil)
	if err != nil {
		return errors.Wrapf(err, "declaring exchange %q", rmq.masterExchangeName())
	}

	// master.delayed 延时Exchange，延时消息发布到该Exchange
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	if err := rmq.channel.ExchangeDeclare(rmq.delayedExchangeName(), "x-delayed-message", true, false, false, false, args); err != nil {
		return errors.Wrapf(err, "declaring exchange %q", rmq.delayedExchangeName())
	}

	// master.retry 重试Exchange，消息处理失败时（3次以内），将消息重新投递给该Exchange
	if err := rmq.channel.ExchangeDeclare(rmq.retryExchangeName(), "direct", true, false, false, false, nil); err != nil {
		return errors.Wrapf(err, "declaring  exchange %q", rmq.retryExchangeName())
	}

	err = declareConsumer(rmq)
	if err != nil {
		return err
	}
	return nil
}

func getEventRouterKey(event Event) string {
	var routerKey string
	if impRouterKey, ok := event.(EventRouterKey); ok {
		routerKey = impRouterKey.RouterKey()
	} else {
		routerKey = utils.StrPoint(utils.StructToName(event))
	}
	return routerKey
}

// declareConsumer declares all queues and bindings for the consumer
func declareConsumer(rmq *RabbitMQ) (err error) {
	for _, listener := range rmq.listeners {
		channel, err := rmq.conn.Channel()
		if err != nil {
			log.Errorf("Failed to open a channel,Error:%s\n", err)
			return err
		}

		queueName := utils.StructToName(listener)
		routerKey := getEventRouterKey(listener.Event())
		if _, err = channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
			return err
		}
		// 判断是否为延时任务
		if listener.Event().Delay() > 0 {
			if err = channel.QueueBind(queueName, routerKey, rmq.delayedExchangeName(), false, nil); err != nil {
				log.Errorf("query bind err", err, rmq.delayedExchangeName())
				return err
			}
		} else {
			if err = channel.QueueBind(queueName, routerKey, rmq.masterExchangeName(), false, nil); err != nil {
				log.Errorf("query bind err", err, rmq.masterExchangeName())
			}
		}

		// 获取消费通道
		err = channel.Qos(1, 0, false)
		if err != nil {
			return err
		}

		published, err := channel.Consume(
			queueName,
			"",
			false,
			false,
			false,
			false,
			nil)
		if err != nil {
			return err
		}

		go func(l Listener, ds <-chan amqp.Delivery) {
			err = rmq.listen(l, ds)
		}(listener, published)
	}

	return
}

// 监听
func (rmq *RabbitMQ) listen(listener Listener, ds <-chan amqp.Delivery) (err error) {
	for msg := range ds {
		retryNums, ok := msg.Headers["retry_nums"].(int32)
		if !ok {
			retryNums = 0
		}
		if err := listener.Handler(msg.Body); err != nil {

			//消息处理失败 进入延时尝试机制
			log.Infof("retry num:", retryNums)
			if retryNums < 3 {
				if err = rmq.retry(listener, msg, retryNums); err != nil {
					log.Errorf("消息处理失败 进入延时尝试机制异常:%s \n", err)
					return err
				}
			} else {
				// 全部处理失败，入库
				log.Warnf("尝试次数耗尽，任务处理失败，触发监听者异常处理函数")

				listener.OnError(msg.Body, err)
				if err = msg.Ack(false); err != nil {
					log.Errorf("确认消息完成异常:%s \n", err)
				}
				return err
			}
			err = msg.Ack(true)
			if err != nil {
				log.Errorf("确认消息未完成异常:%s \n", err)
				return err
			}
		} else {
			err := msg.Ack(false)
			if err != nil {
				log.Errorf("确认消息完成异常:%s \n", err)
				return err
			}
		}
	}
	return nil
}

// Shutdown closes rabbitmq's connection
func (rmq *RabbitMQ) Shutdown() {
	rmq.quitChannel <- true

	log.Info("shutting down rabbitMQ's connection...")

	<-rmq.quitChannel
}

// handleDisconnect handle a disconnection trying to reconnect every 5 seconds
func (rmq *RabbitMQ) handleDisconnect(rs func() error) {
	for {
		select {
		case errChann := <-rmq.closeChannel:
			if errChann != nil {
				log.Errorf("rabbitMQ disconnection: %v", errChann)
			}
		case <-rmq.quitChannel:
			rmq.conn.Close()
			log.Info("...rabbitMQ has been shut down")
			rmq.quitChannel <- true
			return
		}

		log.Info("...trying to reconnect to rabbitMQ...")

		time.Sleep(5 * time.Second)

		if err := rs(); err != nil {
			log.Errorf("rabbitMQ error: %v", err)
		}
	}
}

// Publish sends the given body on the routingKey to the channel
func (rmq *RabbitMQ) Publish(event Event) error {
	return rmq.publish(event)
}

func (rmq *RabbitMQ) publish(event Event) error {
	headers := make(amqp.Table)

	body := event.Body()
	routerKey := getEventRouterKey(event)
	log.Infof("publishing to %q %q", routerKey, event.Body())

	exchange := rmq.masterExchangeName()
	if event.Delay() > 0 {
		headers["x-delay"] = durationToInt(event.Delay())
		exchange = rmq.delayedExchangeName()
	}

	return rmq.channel.Publish(exchange, routerKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         body,
		Headers:      headers,
	})
}

// 消息处理失败之后 延时尝试
func (rmq *RabbitMQ) retry(listener Listener, msg amqp.Delivery, retryNums int32) (err error) {
	//原始路由key
	//原始交换机名
	table := make(amqp.Table)
	table["x-dead-letter-routing-key"] = msg.RoutingKey
	table["x-dead-letter-exchange"] = msg.Exchange
	table["x-message-ttl"] = rmq.retryTTL

	routeKey := fmt.Sprintf("%s.retry", msg.RoutingKey)

	name := utils.StructToName(listener)
	retryQueueName := fmt.Sprintf("%s@retry", name)

	// 定义重试队列
	if _, err = rmq.channel.QueueDeclare(retryQueueName, true, false, false, false, table); err != nil {
		return
	}

	// 队列绑定
	if err = rmq.channel.QueueBind(retryQueueName, routeKey, rmq.retryExchangeName(), false, nil); err != nil {
		return
	}

	// 发送任务
	header := make(amqp.Table)
	header["retry_nums"] = retryNums + 1

	return rmq.channel.Publish(rmq.retryExchangeName(), routeKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         msg.Body,
		Headers:      header,
	})
}
