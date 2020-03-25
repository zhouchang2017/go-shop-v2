package http

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/app/usecases"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"math"
	"net/http"
	"sort"
)

type IndexController struct {
	productSrv *services.ProductService
	topicSrv   *services.TopicService
	articleSrv *services.ArticleService
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
	form.Hidden = "description,attributes,options,images"

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

// topic relation products
// api /topics/:id/products?page=1
func (this *IndexController) TopicProducts(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ResponseError(ctx, err2.Err422.F("缺少id参数"))
		return
	}
	var option request.IndexRequest
	var page int64
	var perPage int64 = 15
	err := ctx.ShouldBind(&option)
	if err == nil {
		page = option.GetPage()
		perPage = option.GetPerPage()
	}
	data, pagination, err := usecases.TopicProductPagination(ctx, id, page, perPage, this.topicSrv, this.productSrv)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	for _, product := range data {
		product.WithMeta("type", product.GetType())
		product.TotalSalesQty += product.FakeSalesQty
		product.FakeSalesQty = 0
	}
	Response(ctx, gin.H{
		"data":       data,
		"pagination": pagination,
	}, http.StatusOK)
}

// product 简约接口，获取 主图，标题，价格，等
func (this *IndexController) Products(ctx *gin.Context) {
	//id := ctx.Query("ids")

}
