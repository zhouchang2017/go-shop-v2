package app

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go-shop-v2/pkg/utils"
	"net/url"
	"time"
)

type Event interface {
	RouterName() string
	Delay() time.Duration
	Body() []byte
}

type Listener interface {
	Event() Event
	OnError(err error)
	Handler(data []byte) error
}

var uri = "amqp://root:12345678@localhost:5672/"

var instance *RabbitMQ

type RabbitMQ struct {
	URL        string
	Exchange   string
	Conn       *amqp.Connection
	Chann      *amqp.Channel
	Queue      amqp.Queue
	closeChann chan *amqp.Error
	quitChann  chan bool
	listeners  []Listener
	retryTTL   int
}

func (r RabbitMQ) masterExchangeName() string {
	return r.Exchange
}

func (r RabbitMQ) delayedExchangeName() string {
	return fmt.Sprintf("%s.delayed", r.Exchange)
}

func (r RabbitMQ) retryExchangeName() string {
	return fmt.Sprintf("%s.retry", r.Exchange)
}

func durationToInt(d time.Duration) int {
	return int(d.Seconds() * 1000)
}

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	VHost    string `json:"vhost"`
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

func New() {

}

// master 主Exchange，发布消息时发布到该Exchange
// master.delayed 延时Exchange，延时消息发布到该Exchange
// master.retry 重试Exchange，消息处理失败时（3次以内），将消息重新投递给该Exchange

func NewSender(uri string, exchange string) (*RabbitMQ, error) {
	duration := time.Second * 5
	rmq := &RabbitMQ{
		URL:      uri,
		Exchange: exchange,
		retryTTL: durationToInt(duration),
	}

	err := rmq.senderLoad()
	if err != nil {
		return nil, err
	}

	rmq.quitChann = make(chan bool)

	go rmq.handleDisconnect(rmq.senderLoad)

	return rmq, err
}

func NewReceive(uri string, exchange string) (*RabbitMQ, error) {
	duration := time.Second * 5
	rmq := &RabbitMQ{
		URL:       uri,
		Exchange:  exchange,
		retryTTL:  durationToInt(duration),
		listeners: listeners,
	}

	err := rmq.load()
	if err != nil {
		return nil, err
	}

	rmq.quitChann = make(chan bool)

	go rmq.handleDisconnect(rmq.load)

	return rmq, err
}

var listeners []Listener

func Register(listener Listener) {
	listeners = append(listeners, listener)
}

func (rmq *RabbitMQ) senderLoad() error {
	var err error

	rmq.Conn, err = amqp.Dial(rmq.URL)
	if err != nil {
		return err
	}

	rmq.Chann, err = rmq.Conn.Channel()
	if err != nil {
		return err
	}

	log.Info("connection to rabbitMQ established")

	rmq.closeChann = make(chan *amqp.Error)
	rmq.Conn.NotifyClose(rmq.closeChann)

	// declare exchange if not exist
	// master 主Exchange，发布消息时发布到该Exchange
	err = rmq.Chann.ExchangeDeclare(rmq.masterExchangeName(), "direct", true, false, false, false, nil)
	if err != nil {
		return errors.Wrapf(err, "declaring exchange %q", rmq.masterExchangeName())
	}

	// master.delayed 延时Exchange，延时消息发布到该Exchange
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	if err := rmq.Chann.ExchangeDeclare(rmq.delayedExchangeName(), "x-delayed-message", true, false, false, false, args); err != nil {
		return errors.Wrapf(err, "declaring exchange %q", rmq.delayedExchangeName())
	}

	// master.retry 重试Exchange，消息处理失败时（3次以内），将消息重新投递给该Exchange
	if err := rmq.Chann.ExchangeDeclare(rmq.retryExchangeName(), "direct", true, false, false, false, nil); err != nil {
		return errors.Wrapf(err, "declaring  exchange %q", rmq.retryExchangeName())
	}
	return nil
}

func (rmq *RabbitMQ) load() error {
	var err error

	rmq.Conn, err = amqp.Dial(rmq.URL)
	if err != nil {
		return err
	}

	rmq.Chann, err = rmq.Conn.Channel()
	if err != nil {
		return err
	}

	log.Info("connection to rabbitMQ established")

	rmq.closeChann = make(chan *amqp.Error)
	rmq.Conn.NotifyClose(rmq.closeChann)

	// declare exchange if not exist
	// master 主Exchange，发布消息时发布到该Exchange
	err = rmq.Chann.ExchangeDeclare(rmq.masterExchangeName(), "direct", true, false, false, false, nil)
	if err != nil {
		return errors.Wrapf(err, "declaring exchange %q", rmq.masterExchangeName())
	}

	// master.delayed 延时Exchange，延时消息发布到该Exchange
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	if err := rmq.Chann.ExchangeDeclare(rmq.delayedExchangeName(), "x-delayed-message", true, false, false, false, args); err != nil {
		return errors.Wrapf(err, "declaring exchange %q", rmq.delayedExchangeName())
	}

	// master.retry 重试Exchange，消息处理失败时（3次以内），将消息重新投递给该Exchange
	if err := rmq.Chann.ExchangeDeclare(rmq.retryExchangeName(), "direct", true, false, false, false, nil); err != nil {
		return errors.Wrapf(err, "declaring  exchange %q", rmq.retryExchangeName())
	}

	err = declareConsumer(rmq)
	if err != nil {
		return err
	}
	return nil
}

// declareConsumer declares all queues and bindings for the consumer
func declareConsumer(rmq *RabbitMQ) (err error) {
	for _, listener := range rmq.listeners {
		channel, err := rmq.Conn.Channel()
		if err != nil {
			log.Errorf("Failed to open a channel,Error:%s\n", err)
			return err
		}

		queueName := utils.StructToName(listener)
		if _, err = channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
			return err
		}
		// 判断是否为延时任务
		if listener.Event().Delay() > 0 {
			if err = channel.QueueBind(queueName, listener.Event().RouterName(), rmq.delayedExchangeName(), false, nil); err != nil {
				log.Errorf("query bind err", err, rmq.delayedExchangeName())
				return err
			}
		} else {
			if err = channel.QueueBind(queueName, listener.Event().RouterName(), rmq.masterExchangeName(), false, nil); err != nil {
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
				log.Warnf("全部处理失败，入库")
				if err = msg.Ack(false); err != nil {
					log.Errorf("确认消息完成异常:%s \n", err)
				}
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
	rmq.quitChann <- true

	log.Info("shutting down rabbitMQ's connection...")

	<-rmq.quitChann
}

// handleDisconnect handle a disconnection trying to reconnect every 5 seconds
func (rmq *RabbitMQ) handleDisconnect(rs func() error) {
	for {
		select {
		case errChann := <-rmq.closeChann:
			if errChann != nil {
				log.Errorf("rabbitMQ disconnection: %v", errChann)
			}
		case <-rmq.quitChann:
			rmq.Conn.Close()
			log.Info("...rabbitMQ has been shut down")
			rmq.quitChann <- true
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
	log.Infof("publishing to %q %q", event.RouterName(), event.Body())

	exchange := rmq.masterExchangeName()
	if event.Delay() > 0 {
		headers["x-delay"] = durationToInt(event.Delay())
		exchange = rmq.delayedExchangeName()
	}

	return rmq.Chann.Publish(exchange, event.RouterName(), false, false, amqp.Publishing{
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

	routeKey := fmt.Sprintf("%s.retry", listener.Event().RouterName())

	name := utils.StructToName(listener)
	retryQueueName := fmt.Sprintf("%s@retry", name)

	// 定义重试队列
	if _, err = rmq.Chann.QueueDeclare(retryQueueName, true, false, false, false, table); err != nil {
		return
	}

	// 队列绑定
	if err = rmq.Chann.QueueBind(retryQueueName, routeKey, rmq.retryExchangeName(), false, nil); err != nil {
		return
	}

	// 发送任务
	header := make(amqp.Table)
	header["retry_nums"] = retryNums + 1

	return rmq.Chann.Publish(rmq.retryExchangeName(), routeKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         msg.Body,
		Headers:      header,
	})
}
