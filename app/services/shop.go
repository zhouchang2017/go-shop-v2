package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func init() {
	register(NewShopService)
}

type ShopService struct {
	rep      *repositories.ShopRep
	adminRep *repositories.AdminRep
}

type shopForm struct {
	Name     string              `json:"name" form:"name" binding:"required"`
	Address  *models.ShopAddress `json:"address"`           // 地址
	Location *models.Location    `json:"location"`          // 坐标
	Members  []string            `json:"members,omitempty"` // 成员
}

// 设置门店成员
func (this *ShopService) SetMembers(ctx context.Context, shop *models.Shop, members ...string) (entity *models.Shop, err error) {
	if len(members) == 0 {
		shop.Members = []*models.AssociatedAdmin{}
	}
	if len(members) > 0 {
		result := <-this.adminRep.FindByIds(ctx, members...)
		if result.Error != nil {
			return nil, result.Error
		}
		for _, admin := range result.Result.([]*models.Admin) {
			shop.Members = append(shop.Members, admin.ToAssociated())
		}
	}
	return shop, nil
}

// 同步门店关联成员
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

func NewShopService(rep *repositories.ShopRep, adminRep *repositories.AdminRep) *ShopService {
	return &ShopService{rep: rep, adminRep: adminRep}
}
