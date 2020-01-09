package charts

import (
	"go-shop-v2/pkg/vue/cards"
)

type Charts struct {
	*cards.Card
	settings map[string]interface{}
	extend   map[string]interface{}
}

func NewCharts() *Charts {
	return &Charts{
		Card:     cards.NewCard(),
		settings: map[string]interface{}{},
		extend:   map[string]interface{}{},
	}
}

func (this *Charts) PrefixComponent() bool {
	return true
}

func (this *Charts) WithSettings(key string, value interface{}) {
	this.settings[key] = value
}

func (this *Charts) SetSettings(settings map[string]interface{}) {
	this.settings = settings
}

func (this *Charts) Settings() map[string]interface{} {
	return this.settings
}

func (this *Charts) Extend() map[string]interface{} {
	return this.extend
}

func (this *Charts) WithExtend(key string, value interface{}) {
	this.extend[key] = value
}

func (this *Charts) SetExtend(extend map[string]interface{}) {
	this.extend = extend
}
