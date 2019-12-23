package resources

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func init() {
	register(NewInventoryResource)
}

// 库存管理
type Inventory struct {
	vue.AbstractResource
	model   *models.Inventory
	rep     *repositories.InventoryRep
	service *services.InventoryService
}

// vue列表页路由创建后置hook
func (i *Inventory) OnIndexRouteCreated(ctx *gin.Context, router *vue.Router) {
	router.WithMeta("Lenses", []map[string]string{
		{
			"name":       "多门店聚合",
			"routerName": fmt.Sprintf("%s.aggregate", vue.StructToSingularLabel(Inventory{})),
		},
	})
}

// 自定义vue路由
func (i *Inventory) CustomVueRouter(ctx *gin.Context, warp *vue.ResourceWarp) []*vue.Router {
	var routers []*vue.Router
	routers = append(routers, i.aggregateRouter(ctx, warp))
	return routers
}

// vue首页聚合路由
func (i *Inventory) aggregateRouter(ctx *gin.Context, warp *vue.ResourceWarp) *vue.Router {
	router := &vue.Router{
		Path:      fmt.Sprintf("%s/aggregate", warp.UriKey()),
		Name:      fmt.Sprintf("%s.aggregate", warp.SingularLabel()),
		Component: fmt.Sprintf(`%s/Aggregate`, warp.UriKey()),
		Hidden:    true,
	}
	router.WithMeta("Title", "多门店聚合")
	router.WithMeta("ResourceName", warp.SingularLabel())
	router.WithMeta("IndexRouterName", warp.IndexRouterName())
	router.WithMeta("IndexTitle", i.Title())
	return router
}

// 自定义api
func (i *Inventory) CustomHttpRouters(router gin.IRouter, uri string, singularLabel string) {
	i.aggregate(router, uri, singularLabel)
}

// 首页聚合
func (i *Inventory) aggregate(router gin.IRouter, uri string, singularLabel string) {
	router.GET(fmt.Sprintf("aggregate/%s", uri, ), func(c *gin.Context) {
		filter := &request.IndexRequest{}
		if err := c.ShouldBind(filter); err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}

		data, pagination, err := i.service.Aggregate(c, filter)
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"pagination": pagination,
			"data":       data,
		})
	})
}

func NewInventoryResource(rep *repositories.InventoryRep, service *services.InventoryService) *Inventory {
	return &Inventory{model: &models.Inventory{}, rep: rep, service: service}
}

func (i *Inventory) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	request.SetSearchField("item.code")
	filters := request.Filters.Unmarshal()
	options := &repositories.QueryOption{}
	err := mapstructure.Decode(filters, options)
	if err != nil {
		return err
	}
	if len(options.Status) > 0 {
		request.AppendFilter("status", bson.D{{"$in", options.Status}})
	}
	if len(options.Shops) > 0 {
		request.AppendFilter("shop.id", bson.D{{"$in", options.Shops}})
	}
	return nil
}

func (i *Inventory) Model() interface{} {
	return i.model
}

func (i *Inventory) Repository() repository.IRepository {
	return i.rep
}

func (i Inventory) Make(model interface{}) vue.Resource {
	return &Inventory{model: model.(*models.Inventory)}
}

func (i *Inventory) SetModel(model interface{}) {
	i.model = model.(*models.Inventory)
}

// 资源主标题
func (i Inventory) Title() string {
	return "库存管理"
}

// 左侧导航栏icon
func (this Inventory) Icon() string {
	return "i-box"
}

// 左侧导航栏分组
func (Inventory) Group() string {
	return "Shop"
}

// 自定义创建Resource按钮文字
func (Inventory) CreateButtonName() string {
	return "产品入库"
}
