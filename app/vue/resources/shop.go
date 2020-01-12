package resources

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/message"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

type Shop struct {
	core.AbstractResource
	model        interface{}
	service      *services.ShopService
	adminService *services.AdminService
}

// 实现列表页api
func (s *Shop) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return s.service.Pagination(ctx, req)
}

// 实现列表页api
func (s *Shop) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return s.service.FindById(ctx, id)
}

// 实现创建api
func (s *Shop) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	form := services.ShopCreateOption{}
	if err := mapstructure.Decode(data, &form); err != nil {
		return "", err
	}

	members := []*models.AssociatedAdmin{}
	if len(form.Members) > 0 {
		admins, err := s.adminService.FindByIds(ctx, form.Members...)
		if err != nil {
			return "", err
		}
		for _, admin := range admins {
			members = append(members, admin.ToAssociated())
		}
	}

	entity, err := s.service.Create(ctx, form, members...)
	// 门店创建事件
	message.Dispatch(events.ShopCreated{Shop: entity})

	return core.CreatedRedirect(s, entity.GetID()), nil
}

func (s *Shop) Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error) {
	form := services.ShopCreateOption{}
	if err := mapstructure.Decode(data, &form); err != nil {
		return "", err
	}

	members := []*models.AssociatedAdmin{}
	if len(form.Members) > 0 {
		admins, err := s.adminService.FindByIds(ctx, form.Members...)
		if err != nil {
			return "", err
		}
		for _, admin := range admins {
			members = append(members, admin.ToAssociated())
		}
	}

	entity, err := s.service.Update(ctx, model.(*models.Shop), form, members...)

	// 门店更新事件
	message.Dispatch(events.ShopUpdated{Shop: entity})

	return core.UpdatedRedirect(s, entity.GetID()), nil
}

func (s *Shop) Make(model interface{}) contracts.Resource {
	return &Shop{
		model:        model,
		service:      s.service,
		adminService: s.adminService,
	}
}

func NewShopResource() *Shop {
	return &Shop{model: &models.Shop{}, service: services.MakeShopService(), adminService: services.MakeAdminService()}
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
				fields.NewTextField("详细地址", "Address.Addr", fields.SetShowOnIndex(false)),
				fields.NewTextField("联系人", "Address.Name"),
				fields.NewTextField("电话", "Address.Phone"),
				fields.NewMapField("位置", "Location"),
			),

			// 更新&创建页面
			fields.NewCheckboxGroup("成员", "Members", fields.OnlyOnForm()).Key("id").CallbackOptions(func() []*fields.CheckboxGroupOption {
				associatedAdmins, _ := s.adminService.AllAssociated(context.Background())
				var adminOptions []*fields.CheckboxGroupOption
				for _, admin := range associatedAdmins {
					adminOptions = append(adminOptions, &fields.CheckboxGroupOption{
						Label: admin.Nickname,
						Value: admin.Id,
					})
				}
				return adminOptions
			}),


			fields.NewTable("成员", "Members", func() []contracts.Field {
				return []contracts.Field{
					fields.NewTextField("ID", "Id", fields.ExceptOnForms()),
					fields.NewTextField("昵称", "Nickname", fields.ExceptOnForms()),
				}
			}),
		}
	}
}

func (s *Shop) Model() interface{} {
	return s.model
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
