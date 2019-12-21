package event

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"sync"
)

type Listener interface {
	Handle(ctx context.Context, event interface{}) error
}

func Dispatch(event interface{}) {
	instance.Dispatch(event)
}

func Listen(listen map[interface{}][]Listener) {
	instance.Listen(listen)
}

func AddListen(event interface{}, listeners ...Listener) {
	instance.AddListen(event, listeners...)
}

func Off(event interface{}) {
	instance.Off(event)
}

var once sync.Once
var instance *Bus

type Bus struct {
	events chan interface{}
	listen map[string][]Listener
}

func NewBus() *Bus {
	once.Do(func() {
		instance = &Bus{
			events: make(chan interface{}, 1024),
			listen: map[string][]Listener{},
		}
	})
	return instance
}

func (b *Bus) Listen(listen map[interface{}][]Listener) {
	for event, listener := range listen {
		eventName := b.eventName(event)
		b.listen[eventName] = listener
	}
}

func (b *Bus) AddListen(event interface{}, listen ...Listener) {
	listens := b.resolverListens(event)
	if len(listens) > 0 {
		listens = append(listens, listen...)
	} else {
		eventName := b.eventName(event)
		b.listen[eventName] = listen
	}
}

func (b *Bus) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Event bus shout down...")
				return
			case event := <-b.events:
				for _, listen := range b.resolverListens(event) {
					go func() {
						select {
						case <-ctx.Done():
							log.Println("Event bus shout down...")
							return
						default:
							log.Printf("Dispatch %s event,Listen on %s\n", b.eventName(event), b.eventName(listen))
							if bytes, err := json.Marshal(event); err == nil {
								log.Printf("Payload %s \n", bytes)
							}
							if err := listen.Handle(ctx, event); err != nil {
								// 记录异常的事件处理
								log.Printf("Dispatch %s event,Listen on %s error:%s\n", b.eventName(event), b.eventName(listen), err)
							}
						}
					}()
				}
			}
		}
	}()
}

func (b *Bus) resolverListens(event interface{}) []Listener {
	eventName := b.eventName(event)
	for name, listeners := range b.listen {
		if eventName == name {
			return listeners
		}
	}
	log.Printf("event[%s] not found!", eventName)
	return []Listener{}
}

func (b *Bus) Dispatch(event interface{}) {
	b.events <- event
}

func (b *Bus) Off(event interface{}) {
	eventName := b.eventName(event)
	for name, _ := range b.listen {
		if eventName == name {
			delete(b.listen, name)
			break
		}
	}
}

func (b *Bus) eventName(i interface{}) string {
	t := reflect.TypeOf(i)
	split := strings.Split(t.String(), ".")
	name := split[len(split)-1]
	return name
}
