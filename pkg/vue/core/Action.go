package core

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/utils"
	"go-shop-v2/pkg/vue/contracts"
)

// 动作URI KEY
func ActionUriKey(action contracts.Action) string {
	if customUri, ok := action.(contracts.CustomUri); ok {
		return customUri.UriKey()
	} else {
		return utils.StructNameToSnakeAndPlural(action)
	}
}

// 列表页动作(批量操作)
func resolverIndexActions(ctx *gin.Context, resource contracts.Resource) []contracts.Action {
	res := []contracts.Action{}

	for _, action := range resource.Actions(ctx) {
		if action.ShowOnIndex() && action.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			res = append(res, action)
		}
	}

	return res
}

//  详情页动作
func resolverDetailActions(ctx *gin.Context, resource contracts.Resource) []contracts.Action {
	res := []contracts.Action{}

	for _, action := range resource.Actions(ctx) {
		if action.ShowOnDetail() && action.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) && action.CanRun(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), resource.Model()) {
			res = append(res, action)
		}
	}

	return res
}

func serializeDetailActions(ctx *gin.Context, resource contracts.Resource) []interface{} {
	data := []interface{}{}
	for _, action := range resolverDetailActions(ctx, resource) {
		data = append(data, serializeAction(ctx, action))
	}
	return data
}

// json
func serializeAction(ctx *gin.Context, action contracts.Action) map[string]interface{} {
	res := make(map[string]interface{})
	res["name"] = action.Name()
	res["uri_key"] = ActionUriKey(action)
	res["component"] = action.Component()
	res["prefix_component"] = action.PrefixComponent()
	fields := []contracts.Field{}
	for _, field := range action.Fields(ctx) {
		if field.AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			fields = append(fields, field)
		}
	}
	res["fields"] = fields
	res["confirm_text"] = action.ConfirmText()
	res["confirm_button_text"] = action.ConfirmButtonText()
	res["cancel_button_text"] = action.CancelButtonText()
	for key, value := range action.Meta() {
		res[key] = value
	}
	return res
}
