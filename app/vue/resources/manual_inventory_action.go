package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/app/vue/pages"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

func init() {
	register(NewInventoryActionResource)
}

type InventoryAction struct {
	core.AbstractResource
	model   interface{}
	service *services.ManualInventoryActionService
}


func (m *InventoryAction) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return m.service.Pagination(ctx,req)
}

func (m *InventoryAction) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return m.service.FindById(ctx,id)
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
			fields.NewTextField("类型", "Type"),
			fields.NewTextField("门店", "Shop.Name"),
			fields.NewTextField("操作者", "User.Nickname"),
			fields.NewTextField("状态", "Status"),
			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),

			panels.NewPanel("商品列表",fields.NewTable("商品列表","Items", func() []contracts.Field {
				return []contracts.Field{
					fields.NewTextField("商品ID","Id"),
					fields.NewTextField("商品货号","Code"),
					fields.NewTextField("品牌","Product.Brand.Name"),
					fields.NewTextField("类目","Product.Category.Name"),
					fields.NewTextField("数量","Qty"),
					fields.NewTextField("状态","Status"),
				}
			})).SetWithoutPending(true),
		}
	}
}

func NewInventoryActionResource(rep *repositories.ManualInventoryActionRep, service *services.ManualInventoryActionService) *InventoryAction {
	return &InventoryAction{
		model:   &models.ManualInventoryAction{},
		service: service,
	}
}

type manualInventoryActionForm struct {
	Type   int8                                `json:"type" form:"type"`
	ShopId string                              `json:"shop_id" form:"shop_id"`
	Items  []*models.ManualInventoryActionItem `json:"items" form:"items"`
	Status int8                                `json:"status" form:"status"`
}

func (m *InventoryAction) UpdateFormParse(ctx *gin.Context, model interface{}) (entity interface{}, err error) {
	form := &manualInventoryActionForm{}
	if err := ctx.ShouldBind(form); err != nil {
		return nil, err
	}
	action := model.(*models.ManualInventoryAction)
	user := ctx2.GetUser(ctx)
	if admin, ok := user.(*models.Admin); ok {
		if err := action.SetType(form.Type); err != nil {
			return nil, err
		}
		// 创建设置为保存状态
		action.SetStatusToSaved()
		action.Items = form.Items
		action.User = admin.ToAssociated()
		if form.ShopId != action.Shop.Id {
			return m.service.SetShop(ctx, action, form.ShopId)
		}
		return action, nil
	}
	return nil, err2.Err401
}

func (m *InventoryAction) CreateFormParse(ctx *gin.Context) (entity interface{}, err error) {
	form := &manualInventoryActionForm{}
	if err := ctx.ShouldBind(form); err != nil {
		return nil, err
	}
	action := &models.ManualInventoryAction{}
	user := ctx2.GetUser(ctx)
	if admin, ok := user.(*models.Admin); ok {
		if err := action.SetType(form.Type); err != nil {
			return nil, err
		}
		// 创建设置为保存状态
		action.SetStatusToSaved()
		action.Items = form.Items
		action.User = admin.ToAssociated()
		return m.service.SetShop(ctx, action, form.ShopId)
	}
	return nil, err2.Err401
}

func (m *InventoryAction) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	return nil
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

func (InventoryAction) Group() string {
	return "Shop"
}


// 自定义页面
func (i *InventoryAction) Pages() []contracts.Page {
	return []contracts.Page{
		pages.NewManualInventoryCreatePage(),
	}
}