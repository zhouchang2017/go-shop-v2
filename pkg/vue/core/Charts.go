package core

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/vue/contracts"
)

// Card URI KEY
func CardUriKey(card contracts.Card) string {
	if customUri, ok := card.(contracts.CustomUri); ok {
		return customUri.UriKey()
	} else {
		return utils.StructNameToSnakeAndPlural(card)
	}
}


// 列表页显示
func resolverIndexCards(ctx *gin.Context, resource contracts.Resource) []contracts.Card {
	cards := []contracts.Card{}
	for _, card := range resource.Cards(ctx) {
		if card.ShowOnIndex() && card.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			cards = append(cards, card)
		}
	}
	return cards
}

// 详情页显示
func resolverDetailCards(ctx *gin.Context, resource contracts.Resource) []contracts.Card {
	cards := []contracts.Card{}
	for _, card := range resource.Cards(ctx) {
		if card.ShowOnDetail() && card.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			cards = append(cards, card)
		}
	}
	return cards
}

func SerializeCard(ctx *gin.Context, card contracts.Card) map[string]interface{} {
	res := make(map[string]interface{})
	if isCharts, ok := card.(contracts.Charts); ok {
		return serializeCharts(ctx, isCharts)
	}

	return res
}

func serializeCharts(ctx *gin.Context, charts contracts.Charts) map[string]interface{} {
	res := make(map[string]interface{})
	res["name"] = charts.Name()
	res["uri_key"] = CardUriKey(charts)
	res["component"] = charts.Component()
	res["prefix_component"] = charts.PrefixComponent()
	res["settings"] = charts.Settings()
	res["extend"] = charts.Extend()
	res["columns"] = charts.Columns()
	res["width"] = charts.Width()

	for key, value := range charts.Meta() {
		res[key] = value
	}
	return res
}

// 列表页cards 序列化
func serializeIndexCards(ctx *gin.Context, resource contracts.Resource) []interface{} {
	data := []interface{}{}
	for _, card := range resolverIndexCards(ctx, resource) {
		data = append(data, SerializeCard(ctx, card))
	}
	return data
}

// 详情页cards 序列化
func serializeDetailCards(ctx *gin.Context, resource contracts.Resource) []interface{} {
	data := []interface{}{}
	for _, card := range resolverDetailCards(ctx, resource) {
		data = append(data, SerializeCard(ctx, card))
	}
	return data
}