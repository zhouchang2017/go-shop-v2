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

// Dashboard页显示
func resolverDashboardCards(ctx *gin.Context) []contracts.Card {
	cards := make([]contracts.Card, 0)
	for _, card := range instance.cardlist {
		if card.ShowOnIndex() && card.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			cards = append(cards, card)
		}
	}
	return cards
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
	res["grid"] = charts.Grid()
	for key, value := range charts.Meta() {
		res[key] = value
	}
	if linkable, ok := charts.(contracts.MoreLink); ok {
		linkable.Link()
		link := make(map[string]interface{})
		link["name"] = linkable.Link().Name()
		link["params"] = linkable.Link().Params()
		link["query"] = linkable.Link().Query()
		res["router"] = link
	}
	if refresh, ok := charts.(contracts.ChartsRefresh); ok {
		duration := 15000
		if refresh.Refresh() > 0 {
			duration = int(refresh.Refresh().Seconds() * 1000)
		}
		res["duration"] = duration
	}
	return res
}

// Dashboard cards 序列化
func serializeDashboardCards(ctx *gin.Context) []interface{} {
	data := make([]interface{}, 0)
	for _, card := range resolverDashboardCards(ctx) {
		data = append(data, SerializeCard(ctx, card))
	}
	return data
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
