package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func init() {
	register(NewAdminService)
}

type AdminService struct {
	rep     *repositories.AdminRep
	shopRep *repositories.ShopRep
}

func NewAdminService(adminRep *repositories.AdminRep, shopRep *repositories.ShopRep) *AdminService {
	return &AdminService{
		rep:     adminRep,
		shopRep: shopRep,
	}
}

// 创建用户
func (a *AdminService) Create(ctx context.Context, model *models.Admin, ShopIds ...string) (admin *models.Admin, err error) {
	shopsRes := <-a.shopRep.FindByIds(ctx, ShopIds...)
	if shopsRes.Error != nil {
		return model, shopsRes.Error
	}
	model.SetShops(shopsRes.Result.([]*models.Shop))
	created := <-a.rep.Create(ctx, model)
	if created.Error != nil {
		return nil, created.Error
	}
	return created.Result.(*models.Admin), nil
}

// 更新用户
func (a *AdminService) Update(ctx context.Context,model *models.Admin, ShopIds ...string) (admin *models.Admin, err error) {
	shopsRes := <-a.shopRep.FindByIds(ctx, ShopIds...)
	if shopsRes.Error != nil {
		return model, shopsRes.Error
	}
	model.SetShops(shopsRes.Result.([]*models.Shop))
	saved := <-a.rep.Save(ctx, model)
	if saved.Error != nil {
		return nil, saved.Error
	}
	return saved.Result.(*models.Admin), nil
}

// 同步成员关联门店
func (a *AdminService) SyncAssociatedShop(ctx context.Context, shop *models.Shop) error {
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
		if _, err := a.rep.Collection().UpdateMany(ctx, filter, updated); err != nil {
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
			if _, err := a.rep.Collection().UpdateMany(ctx, filter, updated); err != nil {
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
			if _, err := a.rep.Collection().UpdateMany(ctx, filter, updated); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AdminService) GetShops(ctx context.Context, model interface{}) (admin *models.Admin, err error) {
	if id, ok := model.(string); ok {
		result := <-a.rep.FindById(ctx, id)
		if result.Error != nil {
			return nil, result.Error
		}
		admin = result.Result.(*models.Admin)
	}
	if m, ok := model.(*models.Admin); ok {
		admin = m
	}
	//if admin != nil {
	//	shopsRes := <-a.shopRep.FindByIds(ctx, admin.GetShopIds()...)
	//	if shopsRes.Error != nil {
	//		return admin, shopsRes.Error
	//	}
	//	admin.SetShops(shopsRes.Result.([]*models.Shop))
	//}
	return admin, nil
}

// 获取所有门店，关联格式输出
func (a *AdminService) AllShops(ctx context.Context) ([]*models.AssociatedShop, error) {
	indexRequest := &request.IndexRequest{}
	indexRequest.Page = -1
	allRes := <-a.shopRep.Pagination(ctx, indexRequest)

	if allRes.Error != nil {
		return nil, allRes.Error
	}
	var res []*models.AssociatedShop
	for _, shop := range allRes.Result.([]*models.Shop) {
		res = append(res, shop.ToAssociated())
	}
	return res, nil
}
