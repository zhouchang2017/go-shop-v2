package http

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	err2 "go-shop-v2/pkg/err"
	"net/http"
)

type ProductController struct {
	productSrv   *services.ProductService
	promotionSrv *services.PromotionService
}

// 获取产品详情
// api /products/:id
func (this *ProductController) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := this.productSrv.FindByIdWithItems(ctx, id)
	if err != nil {
		// err
		ResponseError(ctx, err)
		return
	}

	var items []map[string]interface{}
	var qty int64
	for _, item := range product.Items {
		items = append(items, map[string]interface{}{
			"id":              item.GetID(),
			"code":            item.Code,
			"price":           item.Price,
			"promotion_price": item.PromotionPrice,
			"option_values":   item.OptionValues,
			"qty":             item.Qty,
			"avatar":          item.GetAvatar(),
			"on_sale":         item.OnSale,
		})
		qty += item.Qty
	}

	productResponse := map[string]interface{}{
		"id":              product.GetID(),
		"name":            product.Name,
		"code":            product.Code,
		"brand":           product.Brand,
		"category":        product.Category,
		"attributes":      product.Attributes,
		"avatar":          product.Avatar,
		"options":         product.Options,
		"items":           items,
		"description":     product.Description,
		"price":           product.Price,
		"promotion_price": product.PromotionPrice,
		"images":          product.Images,
		"total_sales_qty": product.TotalSalesQty + product.FakeSalesQty,
		"on_sale":         product.OnSale,
		"sort":            product.Sort,
		"qty":             qty,
	}

	Response(ctx, productResponse, 200)
}

// 获取产品促销活动
// api /products/:id/promotions
func (this *ProductController) Promotion(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ResponseError(ctx, err2.Err422.F("缺少id参数"))
		return
	}
	promotionItems := this.promotionSrv.FindActivePromotionByProductId(ctx, id)
	Response(ctx, promotionItems, http.StatusOK)
}
