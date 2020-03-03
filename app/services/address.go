package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type AddressService struct {
	rep *repositories.AddressRep
}

func NewAddressService(rep *repositories.AddressRep) *AddressService {
	return &AddressService{rep: rep}
}

// 查询地址
func (this *AddressService) FindById(ctx context.Context, id string) (address *models.UserAddress, err error) {
	results := <-this.rep.FindById(ctx, id)
	if results.Error != nil {
		return nil, results.Error
	}
	return results.Result.(*models.UserAddress), nil
}

// 用户地址列表
func (this *AddressService) Pagination(ctx context.Context, req *request.IndexRequest) (addresses []*models.UserAddress, pagination response.Pagination, err error) {
	results := <-this.rep.Pagination(ctx, req)
	if results.Error != nil {
		err = results.Error
		return
	}
	addresses = results.Result.([]*models.UserAddress)
	pagination = results.Pagination
	return
}

// 用户所有地址(小程序用户地址不做分页)
func (this *AddressService) Index(ctx context.Context, userId string) (addresses []*models.UserAddress, err error) {
	req := &request.IndexRequest{OrderBy: "is_default", OrderDirection: request.Filter_DESC, Page: -1}
	req.AppendFilter("user_id", userId)
	result := <-this.rep.Pagination(ctx, req)
	if result.Error != nil {
		return addresses, result.Error
	}
	addresses = result.Result.([]*models.UserAddress)

	if len(addresses) == 0 {
		addresses = []*models.UserAddress{}
	}
	return addresses, nil
}

type UserAddressCreateOption struct {
	UserId       string  `json:"user_id" form:"-"`
	ContactName  *string `json:"contact_name" form:"contact_name"`
	ContactPhone *string `json:"contact_phone" form:"contact_phone"`
	Province     *string `json:"province"`
	City         *string `json:"city"`
	Areas        *string `json:"areas"`
	Addr         *string `json:"addr"`
	IsDefault    int     `json:"is_default" form:"is_default"`
}

// 地址总数
func (this *AddressService) Count(ctx context.Context, userId string) (count int64) {
	results := <-this.rep.Count(ctx, bson.M{"user_id": userId})
	if results.Error != nil {
		return 0
	}
	return results.Result
}

// 添加地址
func (this *AddressService) Create(ctx context.Context, opt *UserAddressCreateOption) (address *models.UserAddress, err error) {
	if opt.UserId == "" {
		return nil, err2.New(http.StatusUnprocessableEntity, "缺少user_id参数")
	}

	count := this.Count(ctx, opt.UserId)

	if opt.IsDefault == 1 && count > 0 {
		// 其他地址设为 0
		this.rep.SetIsDefaultByUserId(ctx, opt.UserId, false)
	}
	// 如果用户地址不存在的情况下，创建则为默认地址
	if count == 0 {
		opt.IsDefault = 1
	}
	model := &models.UserAddress{
		UserId:       opt.UserId,
		ContactName:  *opt.ContactName,
		ContactPhone: *opt.ContactPhone,
		Province:     *opt.Province,
		City:         *opt.City,
		Areas:        *opt.Areas,
		Addr:         *opt.Addr,
		IsDefault:    opt.IsDefault,
	}
	created := <-this.rep.Create(ctx, model)
	if created.Error != nil {
		return nil, created.Error
	}
	return created.Result.(*models.UserAddress), nil
}

// 更新地址
func (this *AddressService) Update(ctx context.Context, model *models.UserAddress) (address *models.UserAddress, err error) {
	if model.IsDefault == 1 {
		// 其他设置为 0
		this.rep.SetIsDefaultByUserId(ctx, model.UserId, false)
	}
	saved := <-this.rep.Save(ctx, model)
	if saved.Error != nil {
		return nil, saved.Error
	}
	return saved.Result.(*models.UserAddress), nil
}

// 删除地址
func (this *AddressService) Delete(ctx context.Context, id string) (err error) {
	return <-this.rep.Delete(ctx, id)
}
