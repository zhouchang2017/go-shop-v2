package http

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/app/tb"
	"go-shop-v2/app/usecases"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"sort"
)

type IndexController struct {
	productSrv   *services.ProductService
	topicSrv     *services.TopicService
	articleSrv   *services.ArticleService
	inventorySrv *services.InventoryService
}

type IndexMorph interface {
	GetSort() int64
}

type dataSlice []IndexMorph

func (d dataSlice) Len() int {
	return len(d)
}

func (d dataSlice) Less(i, j int) bool {
	return d[i].GetSort() > d[j].GetSort()
}

func (d dataSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// morph index,include product、article、topic
func (this *IndexController) Index(ctx *gin.Context) {
	// 处理函数
	form := &request.IndexRequest{}
	if err := ctx.ShouldBind(form); err != nil {
		// error handle
		spew.Dump(err)
	}

	form.AppendFilter("on_sale", true)

	var data dataSlice

	// 不展现 description,attributes,options
	form.Hidden = "description,attributes,options"
	// 只搜索第一张图片
	form.AppendProjection("images", bson.M{"$slice": 1})
	// products
	products, pagination, err := this.productSrv.Pagination(ctx, form)

	if err != nil {
		// error handle
		spew.Dump(err)
	}

	for _, product := range products {
		product.WithMeta("type", product.GetType())
		product.TotalSalesQty += product.FakeSalesQty
		product.FakeSalesQty = 0
		data = append(data, product)
	}

	count := pagination.Total
	var totalPage int64 = 1
	if count < pagination.PerPage {
		totalPage = 1
	} else {
		totalPage = int64(math.Floor(float64(count / pagination.PerPage)))
	}

	// topics
	topicCount := this.topicSrv.Count(ctx)
	topics, _, err := this.topicSrv.SimplePagination(ctx, form.Page, int64(math.Floor(float64(topicCount/totalPage))))
	if err != nil {
		// err
		spew.Dump(err)
	}

	for _, topic := range topics {
		topic.WithMeta("type", topic.GetType())
		data = append(data, topic)
	}
	// articles
	articleCount := this.articleSrv.Count(ctx)
	articles, _, err := this.articleSrv.SimplePagination(ctx, form.Page, int64(math.Floor(float64(articleCount/totalPage))))
	if err != nil {
		// err
		spew.Dump(err)
	}

	for _, article := range articles {
		article.WithMeta("type", article.GetType())
		data = append(data, article)
	}

	sort.Sort(data)

	Response(ctx, gin.H{
		"data":       data,
		"pagination": pagination,
	}, 200)
}

// article detail
func (this *IndexController) article(ctx *gin.Context) {
	id := ctx.Param("id")
	article, err := this.articleSrv.FindById(ctx, id)
	if err != nil {
		// err
		spew.Dump(err)
	}

	Response(ctx, article, 200)
}

// topic detail
func (this *IndexController) Topic(ctx *gin.Context) {
	id := ctx.Param("id")
	topic, err := this.topicSrv.FindById(ctx, id)
	if err != nil {
		// err
		spew.Dump(err)
	}
	Response(ctx, topic, 200)
}

// product detail
func (this *IndexController) Product(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := usecases.ProductWithStock(ctx, id, this.productSrv, this.inventorySrv)
	if err != nil {
		// err
		spew.Dump(err)
	}

	var items []map[string]interface{}

	for _, item := range product.Items {
		items = append(items, map[string]interface{}{
			"id":            item.GetID(),
			"code":          item.Code,
			"price":         item.Price,
			"option_values": item.OptionValues,
			"qty":           item.Qty,
		})
	}

	productResponse := map[string]interface{}{
		"id":              product.GetID(),
		"name":            product.Name,
		"code":            product.Code,
		"brand":           product.Brand,
		"category":        product.Category,
		"attributes":      product.Attributes,
		"options":         product.Options,
		"items":           items,
		"description":     product.Description,
		"price":           product.Price,
		"images":          product.Images,
		"total_sales_qty": product.TotalSalesQty + product.FakeSalesQty,
		"on_sale":         product.OnSale,
		"sort":            product.Sort,
		"qty":             product.Qty,
	}

	Response(ctx, productResponse, 200)
}

// product 简约接口，获取 主图，标题，价格，等
func (this *IndexController) Products(ctx *gin.Context) {
	//id := ctx.Query("ids")

}

type taobaoResponse struct {
	Data *tbResponseBody `json:"data"`
}

type tbResponseBody struct {
	//ApiStack []tbApiStackItem       `json:"apiStack"`
	Item     tbResponseBodyItem     `json:"item"`
	MockData string                 `json:"mockData,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Props    interface{}            `json:"props"`
	PropsCut string                 `json:"propsCut"`
	SkuBase  skuBase                `json:"skuBase"`
}

type tbApiStackItem map[string]interface{}

type skuBase struct {
	Props []tbSkuBaseItem `json:"props"`
}

type tbSkuBaseItem struct {
	Name   string               `json:"name"`
	Pid    string               `json:"pid"`
	Values []tbSkuBaseItemValue `json:"values"`
}

type tbSkuBaseItemValue struct {
	Name string `json:"name"`
	Vid  string `json:"vid"`
}

type tbResponseBodyItem struct {
	BrandValueId    string   `json:"brandValueId"`
	CartUrl         string   `json:"cartUrl"`
	CategoryId      string   `json:"categoryId"`
	CommentCount    string   `json:"commentCount"`
	H5moduleDescUrl string   `json:"h5moduleDescUrl"`
	Images          []string `json:"images"`
	ItemId          string   `json:"itemId"`
	Title           string   `json:"title"`
}

// test taobao product detail api
func (this *IndexController) TaobaoDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		// err
		err2.ErrorEncoder(ctx, errors.New("id 参数缺少"), ctx.Writer)
		return
	}

	service := &tb.TaobaoSdkService{}
	data, err := service.Detail(id)

	if err != nil {
		// err
		err2.ErrorEncoder(ctx, err, ctx.Writer)
		return
	}

	ctx.JSON(200, data)
}