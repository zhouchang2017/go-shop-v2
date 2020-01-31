package http

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"math"
	"sort"
)

type ProductController struct {
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

// 产品列表
func (this *ProductController) Index(ctx *gin.Context) {
	// 处理函数
	form := &request.IndexRequest{}
	if err := ctx.ShouldBind(form); err != nil {
		// error handle
		spew.Dump(err)
	}

	var data dataSlice

	// products
	products, pagination, err := this.productSrv.Pagination(ctx, form)

	if err != nil {
		// error handle
		spew.Dump(err)
	}

	for _, product := range products {
		product.WithMeta("type",product.GetType())
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
		topic.WithMeta("type",topic.GetType())
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
		article.WithMeta("type",article.GetType())
		data = append(data, article)
	}

	sort.Sort(data)

	Response(ctx, gin.H{
		"data":       data,
		"pagination": pagination,
	}, 200)
}
