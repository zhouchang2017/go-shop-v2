package filters

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/filters"
)

type ShopFilter struct {
	*filters.AbstractSelectFilter
	rep *repositories.ShopRep
}


func NewShopFilter() *ShopFilter {
	return &ShopFilter{
		AbstractSelectFilter: filters.NewAbstractSelectFilter().Multiple(),
		rep:                  repositories.NewShopRep(mongodb.GetConFn()),
	}
}

func (this ShopFilter) Apply(ctx *gin.Context, value interface{}, request *request.IndexRequest) error {
	spew.Dump(value)
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
	res := this.rep.GetAllAssociatedShops(ctx)
	data := []contracts.FilterOption{}
	for _, shop := range res {
		data = append(data, filters.NewSelectOption(shop.Name, shop.Id))
	}
	return data
}
