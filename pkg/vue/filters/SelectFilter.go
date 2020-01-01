package filters

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/vue/contracts"
)

type Option struct {
	OptionLabel string      `json:"label"`
	OptionValue interface{} `json:"value"`
}

func NewSelectOption(label string, value interface{}) *Option {
	return &Option{OptionLabel: label, OptionValue: value}
}

func (this Option) Label() string {
	return this.OptionLabel
}

func (this Option) Value() interface{} {
	return this.OptionValue
}

type AbstractSelectFilter struct {
	meta            map[string]interface{}
	prefixComponent bool
}

func NewAbstractSelectFilter() *AbstractSelectFilter {
	return &AbstractSelectFilter{
		meta: map[string]interface{}{},
	}
}

func (AbstractSelectFilter) Component() string {
	return "select-filter"
}

func (AbstractSelectFilter) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (this *AbstractSelectFilter) WithMeta(key string, value interface{}) {
	this.meta[key] = value
}

func (this AbstractSelectFilter) Meta() map[string]interface{} {
	return this.meta
}

func (this AbstractSelectFilter) PrefixComponent() bool {
	return true
}

func (this *AbstractSelectFilter) Multiple() *AbstractSelectFilter {
	this.WithMeta("multiple", true)
	return this
}

func SerializeMap(ctx *gin.Context, filter contracts.Filter) map[string]interface{} {
	res := make(map[string]interface{})
	res["name"] = filter.Name()
	res["key"] = filter.Key()
	res["component"] = filter.Component()
	res["prefix_component"] = filter.PrefixComponent()
	res["value"] = filter.DefaultValue(ctx)
	res["options"] = filter.Options(ctx)
	for key, value := range filter.Meta() {
		res[key] = value
	}
	return res
}
