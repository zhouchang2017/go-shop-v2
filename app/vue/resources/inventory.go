package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	register(NewInventoryResource)
}

// 库存管理
type Inventory struct {
	model   interface{}
	rep     *repositories.InventoryRep
	service *services.InventoryService
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
			fields.NewTextField("货号", "Item.Code"),
			fields.NewTextField("状态", "Status"),
			fields.NewTextField("库存", "Qty"),

			panels.NewPanel("门店信息",
				fields.NewTextField("门店ID", "Shop.Id", fields.OnlyOnDetail()),
				fields.NewTextField("门店", "Shop.Name", fields.OnlyOnDetail()),
			),

			panels.NewPanel("产品信息",
				fields.NewTextField("产品ID", "Item.Product.Id", fields.OnlyOnDetail()),
				fields.NewTextField("产品货号", "Item.Product.Code", fields.OnlyOnDetail()),
				fields.NewTextField("产品名称", "Item.Product.Name", fields.OnlyOnDetail()),
			),

			panels.NewPanel("商品信息",
				fields.NewTextField("商品ID", "Item.Id", fields.OnlyOnDetail()),
				fields.NewTextField("商品货号", "Item.Code", fields.OnlyOnDetail()),
				fields.NewTable("销售属性", "Item.OptionValues", func() []contracts.Field {
					return []contracts.Field{
						fields.NewTextField("编码", "Code"),
						fields.NewTextField("值", "Value"),
					}
				}),
			),
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

func (i Inventory) Make(model interface{}) contracts.Resource {
	return &Inventory{
		rep:     i.rep,
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

//// 自定义聚合
//func (i Inventory) Lenses() []vue.Lens {
//	return []vue.Lens{
//		lenses.NewInventoryAggregateLens(&Inventory{}, i.service),
//	}
//}

//type manualActionsLink struct {
//}
//
//func (manualActionsLink) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
//	admin := user.(*models.Admin)
//	return len(admin.Shops) > 0
//}
//
//func (manualActionsLink) Title() string {
//	return "库存操作"
//}
//
//func (manualActionsLink) RouterName() string {
//	return vue.NewResourceHelper(&ManualInventoryAction{}).IndexRouterName()
//}
//
//// 自定义link
//func (i Inventory) Links() []vue.Link {
//	return []vue.Link{
//		manualActionsLink{},
//	}
//}
