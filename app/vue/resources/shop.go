package resources

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/event"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

func init() {
	register(NewShopResource)
}

type Shop struct {
	model        interface{}
	rep          *repositories.ShopRep
	service      *services.ShopService
	adminService *services.AdminService
}



// 实现列表页api
func (s *Shop) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	results := <-s.rep.Pagination(ctx, req)
	return results.Result, results.Pagination, results.Error
}

// 实现列表页api
func (s *Shop) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	result := <-s.rep.FindById(ctx, id)
	return result.Result, result.Error
}

// 实现创建api
func (s *Shop) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	form := &shopForm{}
	if err := mapstructure.Decode(data, form); err != nil {
		return "", err
	}
	shop := &models.Shop{
		Name:     form.Name,
		Address:  form.Address,
		Location: form.Location,
	}
	entity, err := s.service.SetMembers(ctx, shop, form.Members...)
	if err != nil {
		return "", err
	}
	created := <-s.rep.Create(ctx, entity)
	if created.Error != nil {
		return "", created.Error
	}
	event.Dispatch(events.ShopCreated{Shop: created.Result.(*models.Shop)})
	return core.CreatedRedirect(s, created.Id), nil
}

func (s *Shop) Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error) {
	spew.Dump(data)
	panic("implement me")
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
		model:        model,
		rep:          s.rep,
		service:      s.service,
		adminService: s.adminService,
	}
}

func NewShopResource(rep *repositories.ShopRep, service *services.ShopService, adminService *services.AdminService) *Shop {
	return &Shop{model: &models.Shop{}, rep: rep, service: service, adminService: adminService}
}

type shopForm struct {
	Name     string              `json:"name"`
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

// 字段设置
func (s *Shop) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("名称", "Name"),
			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),

			panels.NewPanel("地址",
				fields.NewAreaCascader("省/市/区", "Address"),
				fields.NewTextField("详细地址", "Address.Addr",fields.SetShowOnIndex(false)),
				fields.NewTextField("联系人", "Address.Name"),
				fields.NewTextField("电话", "Address.Phone"),
				fields.NewMapField("位置", "Location"),
			),

			fields.NewCheckboxGroup("成员", "Members", fields.OnlyOnForm()).Key("id").CallbackOptions(func() []*fields.CheckboxGroupOption {
				associatedAdmins, _ := s.adminService.AllAdmins(context.Background())
				var adminOptions []*fields.CheckboxGroupOption
				for _, admin := range associatedAdmins {
					adminOptions = append(adminOptions, &fields.CheckboxGroupOption{
						Label: admin.Nickname,
						Value: admin.Id,
					})
				}
				return adminOptions
			}),


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
