package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/event"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

func init() {
	register(NewShopResource)
}

type Shop struct {
	model   interface{}
	rep     *repositories.ShopRep
	service *services.ShopService
}

func (s *Shop) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}

func (s *Shop) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (s *Shop) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (s *Shop) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (s *Shop) Policy() interface{} {
	return nil
}

func (s *Shop) Make(model interface{}) contracts.Resource {
	return &Shop{
		model:   model,
		rep:     s.rep,
		service: s.service,
	}
}

func NewShopResource(rep *repositories.ShopRep, service *services.ShopService) *Shop {
	return &Shop{model: &models.Shop{}, rep: rep, service: service}
}

type shopForm struct {
	Name     string              `json:"name" form:"name" binding:"required"`
	Address  *models.ShopAddress `json:"address"`           // 地址
	Location *models.Location    `json:"location"`          // 坐标
	Members  []string            `json:"members,omitempty"` // 成员
}

// 更新表单处理
func (s *Shop) UpdateFormParse(ctx *gin.Context, model interface{}) (entity interface{}, err error) {
	shop := model.(*models.Shop)
	form := &shopForm{}
	if err = ctx.ShouldBind(form); err != nil {
		return nil, err
	}
	shop.Name = form.Name
	shop.Location = form.Location
	shop.Address = form.Address

	return s.service.SetMembers(ctx, shop, form.Members...)
}

// 创建表单处理
func (s *Shop) CreateFormParse(ctx *gin.Context) (entity interface{}, err error) {
	form := &shopForm{}
	if err = ctx.ShouldBind(form); err != nil {
		return nil, err
	}
	shop := &models.Shop{
		Name:     form.Name,
		Address:  form.Address,
		Location: form.Location,
	}

	return s.service.SetMembers(ctx, shop, form.Members...)
}

// 创建成功钩子
func (s *Shop) Created(ctx *gin.Context, resource interface{}) {
	event.Dispatch(events.ShopCreated{Shop: resource.(*models.Shop)})
}

// 更新成功钩子
func (s *Shop) Updated(ctx *gin.Context, resource interface{}) {
	event.Dispatch(events.ShopUpdated{Shop: resource.(*models.Shop)})
}

// 列表页&详情页展示字段设置
func (s *Shop) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("名称", "Name"),
			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),

			panels.NewPanel("地址",
				fields.NewTextField("省份", "Address.Province", fields.OnlyOnDetail()),
				fields.NewTextField("城市", "Address.City", fields.OnlyOnDetail()),
				fields.NewTextField("区/县", "Address.Areas", fields.OnlyOnDetail()),
				fields.NewTextField("详细地址", "Address.Addr", fields.OnlyOnDetail()),
				fields.NewTextField("联系人", "Address.Name", fields.OnlyOnDetail()),
				fields.NewTextField("电话", "Address.Phone", fields.OnlyOnDetail()),

			),


			panels.NewPanel("成员",
				fields.NewTable("成员", "Members", func() []contracts.Field {
					return []contracts.Field{
						fields.NewTextField("ID", "Id", fields.ExceptOnForms()),
						fields.NewTextField("昵称", "Nickname", fields.ExceptOnForms()),
					}
				})).SetWithoutPending(true),
		}
	}
}

// 列表页搜索处理
func (s *Shop) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	return nil
}

func (s *Shop) Model() interface{} {
	return s.model
}

func (s *Shop) Repository() repository.IRepository {
	return s.rep
}

func (s *Shop) SetModel(model interface{}) {
	s.model = model.(*models.Shop)
}

func (s Shop) Title() string {
	return "门店"
}

func (Shop) Icon() string {
	return "icons-store"
}

func (Shop) Group() string {
	return "Shop"
}
