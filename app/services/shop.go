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
	register(NewShopService)
}

type ShopService struct {
	rep *repositories.ShopRep
}

func (this *ShopService) Pagination(ctx context.Context, req *request.IndexRequest) (shops []*models.Shop, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		return shops, pagination, results.Error
	}
	return results.Result.([]*models.Shop), results.Pagination, nil
}

func (this *ShopService) FindById(ctx context.Context, id string) (shop *models.Shop, err error) {
	byId := <-this.rep.FindById(ctx, id)
	if byId.Error != nil {
		err = byId.Error
		return
	}
	return byId.Result.(*models.Shop), nil
}

func (this *ShopService) FindByIds(ctx context.Context, ids ...string) (shops []*models.Shop, err error) {
	byIds := <-this.rep.FindByIds(ctx, ids...)
	if byIds.Error != nil {
		return nil, byIds.Error
	}
	return byIds.Result.([]*models.Shop), nil
}

func (this *ShopService) All(ctx context.Context) (shops []*models.Shop, err error) {
	all := <-this.rep.FindAll(ctx)
	if all.Error != nil {
		return nil, all.Error
	}
	return all.Result.([]*models.Shop), nil
}

func (this *ShopService) AllAssociatedShops(ctx context.Context) (shops []*models.AssociatedShop, err error) {
	shops2, err := this.All(ctx)
	if err != nil {
		return nil, err
	}
	shops = []*models.AssociatedShop{}
	for _, shop := range shops2 {
		shops = append(shops, shop.ToAssociated())
	}
	return shops, nil
}

type ShopCreateOption struct {
	Name     string              `json:"name" form:"name" binding:"required"`
	Address  *models.ShopAddress `json:"address"`           // 地址
	Location *models.Location    `json:"location"`          // 坐标
	Members  []string            `json:"members,omitempty"` // 成员
}

// 创建门店
func (this *ShopService) Create(ctx context.Context, option ShopCreateOption, admins ...*models.AssociatedAdmin) (entity *models.Shop, err error) {
	shop := models.NewShop()
	shop.Name = option.Name
	shop.Address = option.Address
	shop.Location = option.Location
	shop.Members = admins

	created := <-this.rep.Create(ctx, shop)
	if created.Error != nil {
		return nil, created.Error
	}
	return created.Result.(*models.Shop), nil
}

// 更新门店
func (this *ShopService) Update(ctx context.Context, shop *models.Shop, option ShopCreateOption, admins ...*models.AssociatedAdmin) (entity *models.Shop, err error) {
	shop.Name = option.Name
	shop.Address = option.Address
	shop.Location = option.Location
	shop.Members = admins

	saved := <-this.rep.Save(ctx, shop)
	if saved.Error != nil {
		return nil, saved.Error
	}
	return saved.Result.(*models.Shop), nil
}

// 同步门店关联成员(创建用户、更新用户时调用)
func (this *ShopService) SyncAssociatedMembers(ctx context.Context, admin *models.Admin) error {
	shops := admin.Shops

	var objIds []primitive.ObjectID
	for _, shop := range shops {
		ids, err := primitive.ObjectIDFromHex(shop.Id)
		if err != nil {
			log.Printf("sync associated member shop id = %s ,to object id error:%s\n", shop.Id, err)
			continue
		}
		objIds = append(objIds, ids)
	}

	if len(shops) == 0 {
		// 该用户无门店所属
		filter := bson.M{"members.id": admin.GetID()}
		updated := bson.M{"$pull": bson.M{"members": bson.M{"id": admin.GetID()}}}
		_, err := this.rep.Collection().UpdateMany(ctx, filter, updated)
		return err
	}
	{
		// 成员信息变更
		filter := bson.M{"members.id": admin.GetID()}
		updated := bson.M{"$set": bson.M{"members.$.nickname": admin.Nickname}}
		if _, err := this.rep.Collection().UpdateMany(ctx, filter, updated); err != nil {
			return err
		}
	}
	{
		// 添加新成员
		filter := bson.M{
			"_id":        bson.D{{"$in", objIds}},
			"members.id": bson.D{{"$ne", admin.GetID()}},
		}
		updated := bson.M{"$addToSet": bson.M{"members": admin.ToAssociated()}}
		if _, err := this.rep.Collection().UpdateMany(ctx, filter, updated); err != nil {
			return err
		}
	}
	return nil
}

func NewShopService(rep *repositories.ShopRep) *ShopService {
	return &ShopService{rep: rep}
}
