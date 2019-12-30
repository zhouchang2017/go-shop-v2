package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/app/vue/lenses"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue"
	"go.mongodb.org/mongo-driver/bson"
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
	helper  *vue.ResourceHelper
}

func (this *Inventory) OnIndexRouteCreated(ctx *gin.Context, router *vue.Router) {
	// 库存操作权限
	authorizedToAction := false
	if action, ok := this.Root.ResolveWarp(&ManualInventoryAction{}); ok {
		authorizedToAction = action.AuthorizedToCreate(ctx)
		router.WithMeta("ActionButtonText", "库存操作")
		router.WithMeta("ActionRouterName", action.CreateRouterName())
	}
	router.WithMeta("AuthorizedToAction", authorizedToAction)
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
	return &Inventory{model: &models.Inventory{}, rep: rep, service: service, helper: vue.NewResourceHelper(&Inventory{})}
}

func (i *Inventory) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			vue.NewIDField(),
			vue.NewTextField("门店", "Shop.Name"),
			vue.NewTextField("类目", "Item.Product.Category.Name"),
			vue.NewTextField("品牌", "Item.Product.Brand.Name"),
			vue.NewTextField("货号", "Item.Code"),
			vue.NewTextField("状态", "Status"),
			vue.NewTextField("库存", "Qty"),

			vue.NewPanel("门店信息",
				vue.NewTextField("门店ID","Shop.Id",vue.OnlyOnDetail()),
				vue.NewTextField("门店","Shop.Name",vue.OnlyOnDetail()),
				),

			vue.NewPanel("产品信息",
				vue.NewTextField("产品ID","Item.Product.Id",vue.OnlyOnDetail()),
				vue.NewTextField("产品货号","Item.Product.Code",vue.OnlyOnDetail()),
				vue.NewTextField("产品名称","Item.Product.Name",vue.OnlyOnDetail()),
			),

			vue.NewPanel("商品信息",
				vue.NewTextField("商品ID","Item.Id",vue.OnlyOnDetail()),
				vue.NewTextField("商品货号","Item.Code",vue.OnlyOnDetail()),
				vue.NewTable("销售属性","Item.OptionValues", func() []vue.Field {
					return []vue.Field{
						vue.NewTextField("编码","Code"),
						vue.NewTextField("值","Value"),
					}
				}),
			),
		}
	}
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
func (i Inventory) Lenses() []vue.Lens {
	return []vue.Lens{
		lenses.NewInventoryAggregateLens(&Inventory{}, i.service),
	}
}

type manualActionsLink struct {
}

func (manualActionsLink) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	admin := user.(*models.Admin)
	return len(admin.Shops) > 0
}

func (manualActionsLink) Title() string {
	return "库存操作"
}

func (manualActionsLink) RouterName() string {
	return vue.NewResourceHelper(&ManualInventoryAction{}).IndexRouterName()
}

// 自定义link
func (i Inventory) Links() []vue.Link {
	return []vue.Link{
		manualActionsLink{},
	}
}
