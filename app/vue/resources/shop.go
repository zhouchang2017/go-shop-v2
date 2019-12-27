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
	"go-shop-v2/pkg/vue"
)

func init() {
	register(NewShopResource)
}

type Shop struct {
	vue.AbstractResource
	model   *models.Shop
	rep     *repositories.ShopRep
	service *services.ShopService
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
			vue.NewIDField(),
			vue.NewTextField("名称", "Name"),
			vue.NewDateTime("创建时间", "CreatedAt"),
			vue.NewDateTime("更新时间", "UpdatedAt"),

			vue.NewPanel("地址",
				vue.NewTextField("省份", "Address.Province", vue.OnlyOnDetail()),
				vue.NewTextField("城市", "Address.City", vue.OnlyOnDetail()),
				vue.NewTextField("区/县", "Address.Areas", vue.OnlyOnDetail()),
				vue.NewTextField("详细地址", "Address.Addr", vue.OnlyOnDetail()),
				vue.NewTextField("联系人", "Address.Name", vue.OnlyOnDetail()),
				vue.NewTextField("电话", "Address.Phone", vue.OnlyOnDetail()),

			),

			vue.NewPanel("成员",
				vue.NewTable("成员", "Members", map[string]string{
					"ID": "id",
					"昵称": "nickname",
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

func (s Shop) Make(model interface{}) vue.Resource {
	return &Shop{model: model.(*models.Shop)}
}

func (s *Shop) SetModel(model interface{}) {
	s.model = model.(*models.Shop)
}

func (s Shop) Title() string {
	return "门店"
}

func (Shop) Icon() string {
	return "i-store"
}

func (Shop) Group() string {
	return "Shop"
}
