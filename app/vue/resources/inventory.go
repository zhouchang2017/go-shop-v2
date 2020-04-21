package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/app/vue/charts"
	"go-shop-v2/app/vue/filters"
	"go-shop-v2/app/vue/pages"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
	"go.mongodb.org/mongo-driver/bson"
)

// 库存管理
type Inventory struct {
	core.AbstractResource
	model   interface{}
	service *services.InventoryService
}

func (this *Inventory) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return this.service.Pagination(ctx, req)
}

func (this *Inventory) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return this.service.FindById(ctx, id)
}

func (this *Inventory) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}

func (this *Inventory) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (this *Inventory) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (this *Inventory) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (this *Inventory) Policy() interface{} {
	return nil
}

func (this *Inventory) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("门店", "Shop.Name"),
			fields.NewTextField("类目", "Item.Product.Category.Name"),
			fields.NewTextField("品牌", "Item.Product.Brand.Name"),
			fields.NewAvatar("图片", "Item.Avatar").Rounded(),
			fields.NewTextField("货号", "Item.Code"),
			fields.NewStatusField("状态", "Status").WithOptions([]*fields.StatusOption{
				fields.NewStatusOption("良品", 0).Success(),
				fields.NewStatusOption("不良品", 1).Error(),
			}),
			fields.NewTextField("库存", "Qty"),
			fields.NewTextField("锁定库存", "LockedQty"),

			panels.NewPanel("门店信息",
				fields.NewTextField("门店ID", "Shop.RefundNo", fields.OnlyOnDetail()),
				fields.NewTextField("门店", "Shop.Name", fields.OnlyOnDetail()),
			),

			panels.NewPanel("产品信息",
				fields.NewTextField("产品ID", "Item.Product.RefundNo", fields.OnlyOnDetail()),
				fields.NewTextField("产品货号", "Item.Product.Code", fields.OnlyOnDetail()),
				fields.NewTextField("产品名称", "Item.Product.Name", fields.OnlyOnDetail()),
			),


			panels.NewPanel("商品信息",
				fields.NewTextField("商品ID", "Item.RefundNo", fields.OnlyOnDetail()),
				fields.NewTextField("商品货号", "Item.Code", fields.OnlyOnDetail()),
				fields.NewLabelsFields("属性值", "Item.OptionValues").Label("name"),

			),

			fields.NewHasManyField("日志", &InventoryLog{}),
		}
	}
}

func (this *Inventory) ResourceHttpDelete() bool {
	return false
}

func (this *Inventory) ResourceHttpForceDelete() bool {
	return false
}

func (this *Inventory) ResourceHttpRestore() bool {
	return false
}

func NewInventoryResource() *Inventory {
	return &Inventory{model: &models.Inventory{}, service: services.MakeInventoryService()}
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

func (i Inventory) Make(model interface{}) contracts.Resource {
	return &Inventory{
		service: i.service,
		model:   model,
	}
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
	return "icons-box"
}

// 左侧导航栏分组
func (Inventory) Group() string {
	return "Shop"
}

// 自定义创建Resource按钮文字
func (Inventory) CreateButtonName() string {
	return "产品入库"
}

// 自定义聚合
func (i Inventory) Lenses() []contracts.Lens {
	return []contracts.Lens{
		//lenses.NewInventoryAggregatePage(),
	}
}

// 自定义页面
func (i Inventory) Pages() []contracts.Page {
	return []contracts.Page{
		pages.NewInventoryAggregatePage(),
		pages.NewManualInventoryCreatePage(),
	}
}

// 过滤器定义
func (i *Inventory) Filters(ctx *gin.Context) []contracts.Filter {
	return []contracts.Filter{
		filters.NewShopFilter(),
		filters.NewInventoryStatusFilter(),
		filters.NewUpdatedAtFilter(),
	}
}

func (i *Inventory) Cards(ctx *gin.Context) []contracts.Card {
	return []contracts.Card{
		charts.NewShopsInventoryLine(),
		charts.NewShopsInventoryBar(),
	}
}
