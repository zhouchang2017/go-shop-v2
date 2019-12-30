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
	register(NewAdminResource)
}

type Admin struct {
	vue.AbstractResource
	model   *models.Admin
	rep     *repositories.AdminRep
	service *services.AdminService
}

// 字段
func (a *Admin) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			vue.NewIDField(),
			vue.NewTextField("用户名", "Username", vue.SetRules([]*vue.FieldRule{
				{"required", "缺少用户名"},
			})),
			vue.NewTextField("昵称", "Nickname",vue.SetRules([]*vue.FieldRule{
				{"required", "缺少昵称"},
			})),
			vue.NewSelect("用户类型", "Type",vue.SetRules([]*vue.FieldRule{
				{Rule:"required"},
			})).SetOptions([]*vue.SelectOption{
				{Label: "超级管理员", Value: "root"},
				{Label: "管理员", Value: "admin"},
				{Label: "店长", Value: "manager"},
				{Label: "销售员", Value: "salesman"},
			}),
			vue.NewDateTime("创建时间", "CreatedAt"),
			vue.NewDateTime("更新时间", "UpdatedAt"),
			vue.NewPanel("所属门店",
				vue.NewTable("所属门店", "Shops", func() []vue.Field {
					return []vue.Field{
						vue.NewTextField("ID", "Id"),
						vue.NewTextField("门店", "Name"),
					}
				}),
			).SetWithoutPending(true),
		}
	}
}

func NewAdminResource(rep *repositories.AdminRep, service *services.AdminService) *Admin {
	return &Admin{model: &models.Admin{}, rep: rep, service: service}
}

type adminForm struct {
	Username             string                   `json:"username" form:"username"  binding:"required"`
	Password             string                   `json:"password" form:"password"  binding:"required"`
	PasswordConfirmation string                   `json:"password_confirmation" form:"password_confirmation" binding:"required" binding:"eqfield=Password"`
	Nickname             string                   `json:"nickname" form:"nickname"  binding:"required"`
	Type                 string                   `json:"type" form:"type" binding:"required"`
	Shops                []*models.AssociatedShop `json:"shops" form:"shops"`
}

type adminUpdateForm struct {
	Username string                   `json:"username" form:"username"`
	Password string                   `json:"password" form:"password"`
	Nickname string                   `json:"nickname" form:"nickname"`
	Type     string                   `json:"type" form:"type"`
	Shops    []*models.AssociatedShop `json:"shops" form:"shops"`
}

func (a *Admin) UpdateFormParse(ctx *gin.Context, model interface{}) (entity interface{}, err error) {
	form := &adminUpdateForm{}
	err = ctx.ShouldBind(form)
	if err != nil {
		return nil, err
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
	if form.Type != "" {
		if _, err := admin.SetType(form.Type); err != nil {
			return nil, err
		}
	}
	admin.Shops = []*models.AssociatedShop{} // 初始化空数组
	if len(form.Shops) > 0 {
		admin.Shops = form.Shops
	}
	return admin, nil
}

func (a *Admin) CreateFormParse(ctx *gin.Context) (entity interface{}, err error) {
	form := &adminForm{}
	err = ctx.ShouldBind(form)
	if err != nil {
		return nil, err
	}
	admin := &models.Admin{
		Username: form.Username,
		Nickname: form.Nickname,
	}
	if _, err := admin.SetType(form.Type); err != nil {
		return nil, err
	}

	admin.SetPassword(form.Password)
	admin.Shops = []*models.AssociatedShop{} // 初始化空数组
	if len(form.Shops) > 0 {
		admin.Shops = form.Shops
	}
	return admin, nil
}

// 创建成功钩子
func (a *Admin) Created(ctx *gin.Context, resource interface{}) {
	event.Dispatch(events.AdminCreated{Admin: resource.(*models.Admin)})
}

// 更新成功钩子
func (a *Admin) Updated(ctx *gin.Context, resource interface{}) {
	event.Dispatch(events.AdminUpdated{Admin: resource.(*models.Admin)})
}

func (a *Admin) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	return nil
}

func (a *Admin) Model() interface{} {
	return a.model
}

func (a *Admin) Repository() repository.IRepository {
	return a.rep
}

func (a *Admin) getService() *services.AdminService {
	return a.service
}

func (a Admin) Make(model interface{}) vue.Resource {
	return &Admin{model: model.(*models.Admin)}
}

func (a *Admin) SetModel(model interface{}) {
	a.model = model.(*models.Admin)
}

func (a Admin) Title() string {
	return "用户"
}

func (Admin) Icon() string {
	return "icons-user"
}
