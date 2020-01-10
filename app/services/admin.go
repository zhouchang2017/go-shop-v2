package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func init() {
	register(NewAdminService)
}

type AdminService struct {
	rep *repositories.AdminRep
}

func NewAdminService(adminRep *repositories.AdminRep) *AdminService {
	return &AdminService{
		rep: adminRep,
	}
}

// 表单结构
type AdminCreateOption struct {
	Username string `json:"username" `
	Password string `json:"password"`
	//PasswordConfirmation string                   `json:"password_confirmation" form:"password_confirmation" binding:"required" binding:"eqfield=Password"`
	Nickname string   `json:"nickname" `
	Type     string   `json:"type" `
	Shops    []string `json:"shops" `
}

// 列表
func (this *AdminService) Pagination(ctx context.Context, req *request.IndexRequest) (admins []*models.Admin, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	return results.Result.([]*models.Admin), results.Pagination, nil
}

// 创建用户
func (this *AdminService) Create(ctx context.Context, option AdminCreateOption, shops ...*models.AssociatedShop) (admin *models.Admin, err error) {
	model := &models.Admin{
		Username: option.Username,
		Nickname: option.Nickname,
		Type:     option.Type,
	}

	model.SetPassword(option.Password)


	model.Shops = shops
	created := <-this.rep.Create(ctx, model)
	if created.Error != nil {
		return nil, created.Error
	}

	admin = created.Result.(*models.Admin)

	return
}

// 更新用户
func (this *AdminService) Update(ctx context.Context, model *models.Admin, option AdminCreateOption,shops ...*models.AssociatedShop) (admin *models.Admin, err error) {

	if option.Username != "" {
		model.Username = option.Username
	}
	if option.Nickname != "" {
		model.Nickname = option.Nickname
	}
	if option.Password != "" {
		model.SetPassword(option.Password)
	}
	model.Type = option.Type

	model.Shops = shops
	saved := <-this.rep.Save(ctx, model)
	if saved.Error != nil {
		return nil, saved.Error
	}

	admin = saved.Result.(*models.Admin)
	return
}

// 同步成员关联门店（shop模型更新时调用）
func (this *AdminService) SyncAssociatedShop(ctx context.Context, shop *models.Shop) error {
	admins := shop.Members
	var objIds []primitive.ObjectID
	for _, admin := range admins {
		ids, err := primitive.ObjectIDFromHex(admin.Id)
		if err != nil {
			log.Printf("sync associated shop admin id = %s ,to object id error:%s\n", admin.Id, err)
			continue
		}
		objIds = append(objIds, ids)
	}
	// TODO 开启事务
	// 目前不属于该门店的用户需要移除
	{
		// 包含该门店的成员
		filter := bson.M{
			"shops.id": shop.GetID(),
		}
		if len(objIds) > 0 {
			// 目前该门店不包含的用户
			filter["_id"] = bson.D{{"$nin", objIds}}
		}
		// 删除门店
		updated := bson.M{"$pull": bson.M{"shops": bson.M{"id": shop.GetID()}}}
		if _, err := this.rep.Collection().UpdateMany(ctx, filter, updated); err != nil {
			return err
		}

	}
	// 目前属于该门店的成员需要更新门店信息
	{
		if len(objIds) > 0 {
			// 包含该门店的用户
			filter := bson.M{
				"shops.id": shop.GetID(),
			}
			updated := bson.M{"$set": bson.M{"shops.$.name": shop.Name}}
			if _, err := this.rep.Collection().UpdateMany(ctx, filter, updated); err != nil {
				return err
			}
		}
	}
	// 新门店新成员
	{
		if len(objIds) > 0 {
			// 包含该门店的用户
			filter := bson.M{
				"shops.id": bson.D{{"$ne", shop.GetID()}},
				"_id":      bson.D{{"$in", objIds}},
			}
			updated := bson.M{"$push": bson.M{"shops": shop.ToAssociated()}}
			if _, err := this.rep.Collection().UpdateMany(ctx, filter, updated); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *AdminService) WithShops(ctx context.Context, model interface{}) (admin *models.Admin, err error) {
	if id, ok := model.(string); ok {
		result := <-this.rep.FindById(ctx, id)
		if result.Error != nil {
			return nil, result.Error
		}
		admin = result.Result.(*models.Admin)
	}
	if m, ok := model.(*models.Admin); ok {
		admin = m
	}
	return admin, nil
}

// 通过id获取
func (this *AdminService) FindById(ctx context.Context, id string) (admin *models.Admin, err error) {
	res := <-this.rep.FindById(ctx, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result.(*models.Admin), nil
}

// 通过id集合获取
func (this *AdminService) FindByIds(ctx context.Context, ids ...string) (admins []*models.Admin, err error) {
	byIds := <-this.rep.FindByIds(ctx, ids...)
	if byIds.Error != nil {
		return nil, byIds.Error
	}
	return byIds.Result.([]*models.Admin), nil
}

// 删除
func (this *AdminService) Destroy(ctx context.Context, id string) error {
	return <-this.rep.Delete(ctx, id)
}

// 还原
func (this *AdminService) Restore(ctx context.Context, id string) (admin *models.Admin, err error) {
	restored := <-this.rep.Restore(ctx, id)
	if restored.Error != nil {
		return nil, restored.Error
	}
	return restored.Result.(*models.Admin), nil
}

// 获取所有管理员，关联格式输出
func (this *AdminService) AllAssociated(ctx context.Context) ([]*models.AssociatedAdmin, error) {
	all := <-this.rep.FindAll(ctx)
	if all.Error != nil {
		return nil, all.Error
	}
	res := []*models.AssociatedAdmin{}
	for _, admin := range all.Result.([]*models.Admin) {
		res = append(res, admin.ToAssociated())
	}
	return res, nil
}
