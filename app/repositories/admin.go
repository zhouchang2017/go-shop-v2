package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

func init() {
	register(NewAdminRep)
}

type AdminRep struct {
	*mongoRep
}

//func (this *AdminRep) FindById(ctx context.Context, id string) (admin *models.Admin, err error) {
//	byId := <-this.mongoRep.FindById(ctx, id)
//	if byId.Error != nil {
//		return nil, byId.Error
//	}
//	return byId.Result.(*models.Admin), nil
//}
//
//func (this *AdminRep) FindByIds(ctx context.Context, ids ...string) (admins []*models.Admin, err error) {
//	byIds := <-this.mongoRep.FindByIds(ctx, ids...)
//	if byIds.Error != nil {
//		return nil, byIds.Error
//	}
//	return byIds.Result.([]*models.Admin), nil
//}
//
//func (this *AdminRep) FindOne(ctx context.Context, credentials map[string]interface{}) (admin *models.Admin, err error) {
//	one := <-this.mongoRep.FindOne(ctx, credentials)
//	if one.Error != nil {
//		return nil, one.Error
//	}
//	return one.Result.(*models.Admin), nil
//}
//
//func (this *AdminRep) FindAll(ctx context.Context) (admins []*models.Admin, err error) {
//	all := <-this.mongoRep.FindAll(ctx)
//	if all.Error != nil {
//		return nil, all.Error
//	}
//	return all.Result.([]*models.Admin), nil
//}
//
//func (this *AdminRep) Pagination(ctx context.Context, req *request.IndexRequest) (admins []*models.Admin, pagination response.Pagination, err error) {
//	results := <-this.mongoRep.Pagination(ctx, req)
//	if results.Error != nil {
//		err = results.Error
//		return
//	}
//	return results.Result.([]*models.Admin), results.Pagination, results.Error
//}
//
//func (this *AdminRep) Create(ctx context.Context, entity *models.Admin) (admin *models.Admin, err error) {
//
//	created := <-this.mongoRep.Create(ctx, entity)
//
//	if created.Error != nil {
//		return nil, created.Error
//	}
//
//	return created.Result.(*models.Admin), nil
//}
//
//func (this *AdminRep) Save(ctx context.Context, entity *models.Admin) (admin *models.Admin, err error) {
//	saved := <-this.mongoRep.Save(ctx, entity)
//
//	if saved.Error != nil {
//		return nil, saved.Error
//	}
//
//	return saved.Result.(*models.Admin), nil
//}
//
//func (this *AdminRep) Update(ctx context.Context, id string, update interface{}) (admin *models.Admin, err error) {
//	results := <-this.mongoRep.Update(ctx, id, update)
//
//	if results.Error != nil {
//		return nil, results.Error
//	}
//
//	return results.Result.(*models.Admin), nil
//}
//
//func (this *AdminRep) Delete(ctx context.Context, id string) (err error) {
//	return <-this.mongoRep.Delete(ctx, id)
//}
//
//func (this *AdminRep) Restore(ctx context.Context, id string) (admin *models.Admin, err error) {
//	restore := <-this.mongoRep.Restore(ctx, id)
//
//	if restore.Error != nil {
//		return nil, restore.Error
//	}
//
//	return restore.Result.(*models.Admin), nil
//}

func NewAdminRep(con *mongodb.Connection) *AdminRep {
	return &AdminRep{
		mongoRep: NewBasicMongoRepositoryByDefault(&models.Admin{}, con),
	}
}
