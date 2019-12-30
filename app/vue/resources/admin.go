package resources

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/event"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

func init() {
	register(NewAdminResource)
}

type Admin struct {
	model   interface{}
	rep     *repositories.AdminRep
	service *services.AdminService
}



// 实现列表页
func (a *Admin) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	results := <-a.rep.Pagination(ctx, req)
	return results.Result, results.Pagination, results.Error
}

// 实现详情页
func (a *Admin) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	result := <-a.rep.FindById(ctx, id)
	return result.Result, result.Error
}

type adminCreateForm struct {
	Username string `json:"username" `
	Password string `json:"password"`
	//PasswordConfirmation string                   `json:"password_confirmation" form:"password_confirmation" binding:"required" binding:"eqfield=Password"`
	Nickname string   `json:"nickname" `
	Type     string   `json:"type" `
	Shops    []string `json:"shops" `
}

// 实现创建
func (a *Admin) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	form := &adminCreateForm{}
	if err := mapstructure.Decode(data, form); err != nil {
		return "", err
	}
	model := &models.Admin{
		Username: form.Username,
		Nickname: form.Nickname,
		Type:     form.Type,
	}

	model.SetPassword(form.Password)

	admin, err := a.service.Create(ctx, model, form.Shops...)

	if err != nil {
		return "", err
	}

	event.Dispatch(events.AdminCreated{Admin: admin})

	return core.CreatedRedirect(a, admin.GetID()), nil
}

// 实现更新
func (a *Admin) Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error) {
	form := &adminCreateForm{}
	if err := mapstructure.Decode(data, form); err != nil {
		return "", err
	}

	admin := model.(*models.Admin)
	if form.Username != "" {
		admin.Username = form.Username
	}
	if form.Nickname != "" {
		admin.Nickname = form.Nickname
	}
	if form.Password != "" {
		admin.SetPassword(form.Password)
	}
	admin.Type = form.Type

	admin2, err := a.service.Update(ctx, admin, form.Shops...)
	if err != nil {
		return "", err
	}

	event.Dispatch(events.AdminUpdated{Admin: admin2})

	return core.UpdatedRedirect(a, admin2.GetID()), nil
}

// 实现删除
func (a *Admin) Destroy(ctx *gin.Context, id string) (err error) {
	return <- a.rep.Delete(ctx, id)
}

func (a Admin) Group() string {
	return "App"
}

func (a *Admin) DisplayInNavigation(ctx *gin.Context, user interface{}) bool {
	return true
}

func (a *Admin) HasIndexRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (a *Admin) HasDetailRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (a *Admin) HasEditRoute(ctx *gin.Context, user interface{}) bool {
	return true
}

func (a *Admin) Policy() interface{} {
	return nil
}

// 字段
func (a *Admin) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {

		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("用户名", "Username", fields.SetRules([]*fields.FieldRule{
				{Rule: "required"},
			})),
			fields.NewTextField("昵称", "Nickname", fields.SetRules([]*fields.FieldRule{
				{Rule: "required"},
			})),

			fields.NewSelect("用户类型", "Type", fields.SetRules([]*fields.FieldRule{
				{Rule: "required"},
			})).SetOptions([]contracts.Field{
				fields.NewTextField("超级管理员", "root", fields.SetValue("root")),
				fields.NewTextField("管理员", "admin", fields.SetValue("admin")),
				fields.NewTextField("店长", "manager", fields.SetValue("manager")),
				fields.NewTextField("销售员", "salesman", fields.SetValue("salesman")),
			}),

			fields.NewPasswordField("密码", "Password", fields.SetRules([]*fields.FieldRule{
				{Rule: "min:6"},
				{Rule: "max:20"},
			}), fields.SetShowOnIndex(false)),

			fields.NewCheckboxGroup("所属门店", "Shops", fields.OnlyOnForm()).Key("id").CallbackOptions(func() []*fields.CheckboxGroupOption {
				associatedShops, _ := a.service.AllShops(context.Background())
				var shopOptions []*fields.CheckboxGroupOption
				for _, shop := range associatedShops {
					shopOptions = append(shopOptions, &fields.CheckboxGroupOption{
						Label: shop.Name,
						Value: shop.Id,
					})
				}
				return shopOptions
			}),

			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),
			panels.NewPanel("所属门店",
				fields.NewTable("所属门店", "Shops", func() []contracts.Field {
					return []contracts.Field{
						fields.NewTextField("ID", "Id"),
						fields.NewTextField("门店", "Name"),
					}
				}),
			).SetWithoutPending(true),
		}
	}
}

func NewAdminResource(rep *repositories.AdminRep, service *services.AdminService) *Admin {
	return &Admin{model: &models.Admin{}, rep: rep, service: service}
}


func (a *Admin) Model() interface{} {
	return a.model
}


func (a Admin) Make(model interface{}) contracts.Resource {
	return &Admin{
		rep:     a.rep,
		service: a.service,
		model:   model,
	}
}

func (a *Admin) SetModel(model interface{}) {
	a.model = model
}

func (a Admin) Title() string {
	return "用户"
}

func (Admin) Icon() string {
	return "icons-user"
}
