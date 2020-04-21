package resources

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/listeners"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go.mongodb.org/mongo-driver/mongo"
)

type Admin struct {
	core.AbstractResource
	model       interface{}
	service     *services.AdminService
	shopService *services.ShopService
}

// 实现列表页
func (a *Admin) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	return a.service.Pagination(ctx, req)
}

// 实现详情页
func (a *Admin) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return a.service.FindById(ctx, id)
}

// 实现创建
func (a *Admin) Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error) {
	option := services.AdminCreateOption{}
	if err := mapstructure.Decode(data, &option); err != nil {
		return "", err
	}

	var admin *models.Admin

	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return "", err
	}
	if err = session.StartTransaction(); err != nil {
		return "", err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// 查找管理门店
		associatedShops := make([]*models.AssociatedShop, 0)
		if len(option.Shops) > 0 {
			shops, err := a.shopService.FindByIds(sessionContext, option.Shops...)
			if err != nil {
				session.AbortTransaction(sessionContext)
				return err
			}
			for _, shop := range shops {
				associatedShops = append(associatedShops, shop.ToAssociated())
			}
		}

		admin, err = a.service.Create(sessionContext, option, associatedShops...)

		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}

		if err := a.shopService.SyncAssociatedMembers(sessionContext, admin); err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}
		session.CommitTransaction(sessionContext)
		return nil
	})
	session.EndSession(ctx)
	if err != nil {
		return "", err
	}

	return core.CreatedRedirect(a, admin.GetID()), nil
}

// 实现更新
func (a *Admin) Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error) {
	option := services.AdminCreateOption{}
	if err := mapstructure.Decode(data, &option); err != nil {
		return "", err
	}

	var admin *models.Admin
	admin = model.(*models.Admin)

	// 开启事务
	var session mongo.Session
	if session, err = mongodb.GetConFn().Client().StartSession(); err != nil {
		return "", err
	}
	if err = session.StartTransaction(); err != nil {
		return "", err
	}
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// 查找管理门店
		associatedShops := make([]*models.AssociatedShop, 0)
		if len(option.Shops) > 0 {
			shops, err := a.shopService.FindByIds(sessionContext, option.Shops...)
			if err != nil {
				session.AbortTransaction(sessionContext)
				return err
			}
			for _, shop := range shops {
				associatedShops = append(associatedShops, shop.ToAssociated())
			}
		}

		admin, err = a.service.Update(sessionContext, admin, option, associatedShops...)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}

		if err := a.shopService.SyncAssociatedMembers(sessionContext, admin); err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}
		session.CommitTransaction(sessionContext)
		return nil
	})

	session.EndSession(ctx)
	if err != nil {
		return "", err
	}

	return core.UpdatedRedirect(a, admin.GetID()), nil
}

// 实现删除
func (a *Admin) Destroy(ctx *gin.Context, id string) (err error) {
	return a.service.Destroy(ctx, id)
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

			fields.NewSelectField("用户类型", "Type", fields.SetRules([]*fields.FieldRule{
				{Rule: "required"},
			})).WithOptions([]*fields.SelectOption{
				{Label: "超级管理员", Value: "root"},
				{Label: "管理员", Value: "admin"},
				{Label: "店长", Value: "manager"},
				{Label: "销售员", Value: "salesman"},
			}),

			fields.NewPasswordField("密码", "Password", fields.SetRules([]*fields.FieldRule{
				{Rule: "min:6"},
				{Rule: "max:20"},
			}), fields.SetShowOnIndex(false)),

			fields.NewTextField("Email", "Email", fields.SetRules([]*fields.FieldRule{
				{Rule: "required"},
				{Rule: "email"},
			})).Email(),

			fields.NewCheckboxGroup("所属门店", "Shops", fields.OnlyOnForm()).Key("id").CallbackOptions(func() []*fields.CheckboxGroupOption {
				associatedShops, _ := a.shopService.AllAssociatedShops(context.Background())
				var shopOptions []*fields.CheckboxGroupOption
				for _, shop := range associatedShops {
					shopOptions = append(shopOptions, &fields.CheckboxGroupOption{
						Label: shop.Name,
						Value: shop.Id,
					})
				}
				return shopOptions
			}),
			fields.NewCheckboxGroup("通知设置", "Notifies").CallbackOptions(func() []*fields.CheckboxGroupOption {
				var options []*fields.CheckboxGroupOption
				for _, notify := range listeners.GetAppNotifications() {
					options = append(options, &fields.CheckboxGroupOption{
						Label: notify["name"],
						Value: notify["key"],
					})
				}
				return options
			}),

			fields.NewDateTime("创建时间", "CreatedAt"),
			fields.NewDateTime("更新时间", "UpdatedAt"),
			fields.NewTable("所属门店", "Shops", func() []contracts.Field {
				return []contracts.Field{
					fields.NewTextField("ID", "RefundNo"),
					fields.NewTextField("门店", "Name"),
				}
			}),
		}
	}
}

func NewAdminResource() *Admin {
	return &Admin{model: &models.Admin{}, service: services.MakeAdminService(), shopService: services.MakeShopService()}
}

func (a *Admin) Model() interface{} {
	return a.model
}

func (a Admin) Make(model interface{}) contracts.Resource {
	return &Admin{
		service:     a.service,
		shopService: a.shopService,
		model:       model,
	}
}

func (a *Admin) SetModel(model interface{}) {
	a.model = model
}

func (a Admin) Title() string {
	return "管理员"
}

func (Admin) Icon() string {
	return "icons-user"
}
