package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/app/vue/actions"
	"go-shop-v2/app/vue/pages"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
)

type InventoryAction struct {
	core.AbstractResource
	model   interface{}
	service *services.ManualInventoryActionService
}

func (m *InventoryAction) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return m.service.Pagination(ctx, req)
}

func (m *InventoryAction) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return m.service.FindById(ctx, id)
}

// 自定义创建页
func (this *InventoryAction) CreationComponent() contracts.Page {
	return pages.NewManualInventoryCreatePage()
}

// 自定义更新页
func (this *InventoryAction) UpdateComponent() contracts.Page {
	return pages.NewManualInventoryUpdatePage()
}

func (m *InventoryAction) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}

func (m *InventoryAction) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (m *InventoryAction) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (m *InventoryAction) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (m *InventoryAction) Policy() interface{} {
	return nil
}

func (m *InventoryAction) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewStatusField("类型", "Type").WithOptions([]*fields.StatusOption{
				fields.NewStatusOption("入库", 0),
				fields.NewStatusOption("出库", 1),
			}),
			fields.NewTextField("门店", "Shop.Name"),
			fields.NewTextField("操作者", "User.Nickname"),
			fields.NewStatusField("状态", "Status").WithOptions([]*fields.StatusOption{
				fields.NewStatusOption("未提交", 0).Warning(),
				fields.NewStatusOption("完成", 1).Success(),
				fields.NewStatusOption("取消", 2).Cancel(),
			}),
			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),

			fields.NewTable("商品列表", "Items", func() []contracts.Field {
				return []contracts.Field{
					//fields.NewTextField("商品ID", "RefundNo"),
					fields.NewTextField("商品货号", "Code"),
					fields.NewAvatar("图片","Avatar").Rounded(),
					fields.NewTextField("品牌", "Product.Brand.Name"),
					fields.NewTextField("类目", "Product.Category.Name"),
					fields.NewTextField("数量", "Qty"),
					fields.NewStatusField("状态", "Status").WithOptions([]*fields.StatusOption{
						fields.NewStatusOption("良品", 0).Success(),
						fields.NewStatusOption("不良品", 1).Error(),
					}),
				}
			}),
		}
	}
}

func NewInventoryActionResource() *InventoryAction {
	return &InventoryAction{
		model:   &models.ManualInventoryAction{},
		service: services.MakeManualInventoryActionService(),
	}
}

func (m *InventoryAction) Model() interface{} {
	return m.model
}

func (m InventoryAction) Make(model interface{}) contracts.Resource {
	return &InventoryAction{
		model:   model,
		service: m.service,
	}
}

func (m *InventoryAction) SetModel(model interface{}) {
	m.model = model.(*models.ManualInventoryAction)
}

func (m InventoryAction) Title() string {
	return "库存操作"
}

// 左侧导航栏icon
func (this InventoryAction) Icon() string {
	return "icons-flag"
}

func (InventoryAction) Group() string {
	return "Shop"
}

// 自定义页面
func (i *InventoryAction) Pages() []contracts.Page {
	return []contracts.Page{
		pages.NewManualInventoryCreatePage(),
	}
}

// 动作
func (i *InventoryAction) Actions(ctx *gin.Context) []contracts.Action {
	return []contracts.Action{
		actions.NewInventoryActionToFinished(),
	}
}
