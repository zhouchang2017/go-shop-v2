package core

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/filters"
)

func serializeFilters(ctx *gin.Context, f ...contracts.Filter) []map[string]interface{} {
	data := []map[string]interface{}{}
	for _, filter := range f {
		data = append(data, filters.Serialize(ctx, filter))
	}
	return data
}

// 资源过滤器
func resolverFilters(target interface{ Filters(ctx *gin.Context) []contracts.Filter }, c *gin.Context) []contracts.Filter {
	filters := []contracts.Filter{}
	for _, filter := range target.Filters(c) {
		// 验证权限
		if filter.AuthorizedTo(c, ctx.GetUser(c).(auth.Authenticatable)) {
			filters = append(filters, filter)
		}
	}
	return filters
}
