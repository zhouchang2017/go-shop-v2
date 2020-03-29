package app

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"go-shop-v2/pkg/utils"
	"time"
)

type OrderTimeOutEvent struct {
}

func (o OrderTimeOutEvent) RouterName() string {
	return "order.timeout"
}

func (o OrderTimeOutEvent) Delay() time.Duration {
	return time.Second * 3
}

func (o OrderTimeOutEvent) Body() []byte {
	return []byte(utils.RandomOrderNo("超时订单"))
}

type OrderCreated struct {
}

func (o OrderCreated) RouterName() string {
	return "order.created"
}

func (o OrderCreated) Delay() time.Duration {
	return time.Second * 0
}

func (o OrderCreated) Body() []byte {
	return []byte(utils.RandomOrderNo("新订单订单号"))
}

type OnOrderCreatedNotifyAdmin struct {
}

func (o OnOrderCreatedNotifyAdmin) Event() Event {
	return OrderCreated{}
}

func (o OnOrderCreatedNotifyAdmin) OnError(err error) {
	panic("implement me")
}

func (o OnOrderCreatedNotifyAdmin) Handler(data []byte) error {
	log.Infof("订单创建事件，后台admin通知处理", string(data))
	return errors.New("创建事件处理失败")
}

type OnOrderTimeoutNotifyAdmin struct {
}

func (o OnOrderTimeoutNotifyAdmin) Event() Event {
	return OrderTimeOutEvent{}
}

func (o OnOrderTimeoutNotifyAdmin) OnError(err error) {
	panic("implement me")
}

func (o OnOrderTimeoutNotifyAdmin) Handler(data []byte) error {
	log.Infof("订单超时事件，后台admin通知处理", string(data))
	return errors.New("后台admin通知处理失败")
}

type OnOrderTimoutNotifyUser struct {
}

func (o OnOrderTimoutNotifyUser) Event() Event {
	return OrderTimeOutEvent{}
}

func (o OnOrderTimoutNotifyUser) OnError(err error) {
	panic("implement me")
}

func (o OnOrderTimoutNotifyUser) Handler(data []byte) error {
	log.Infof("订单超时事件，前台user通知处理", string(data))
	return errors.New("前台user通知处理失败")
}
