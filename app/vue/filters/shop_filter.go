package filters

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/filters"
)

type ShopFilter struct {
	*filters.AbstractSelectFilter
	service *services.ShopService
}

func NewShopFilter() *ShopFilter {
	return &ShopFilter{
		AbstractSelectFilter: filters.NewAbstractSelectFilter().Multiple(),
		service:              services.MakeShopService(),
	}
}

func (this ShopFilter) Apply(ctx *gin.Context, value interface{}, request *request.IndexRequest) error {
	return nil
}

func (this ShopFilter) Key() string {
	return "shops"
}

func (this ShopFilter) Name() string {
	return "门店"
}

func (this ShopFilter) DefaultValue(ctx *gin.Context) interface{} {
	return []interface{}{}
}

func (this ShopFilter) Options(ctx *gin.Context) []contracts.FilterOption {
	shops, _ := this.service.AllAssociatedShops(ctx)
	data := []contracts.FilterOption{}
	for _, shop := range shops {
		data = append(data, filters.NewSelectOption(shop.Name, shop.Id))
	}
	return data
}
