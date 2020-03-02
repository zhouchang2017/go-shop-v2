package services

import (
	"context"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type AddressService struct {
	rep *repositories.AddressRep
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

type UserAddressCreateOption struct {
	UserId       string `json:"user_id" form:"-"`
	ContactName  string `json:"contact_name" form:"contact_name"`
	ContactPhone string `json:"contact_phone" form:"contact_phone"`
	Province     string `json:"province"`
	City         string `json:"city"`
	Areas        string `json:"areas"`
	Addr         string `json:"addr"`
	IsDefault    int    `json:"is_default" form:"is_default"`
}

// 添加地址
func (this *AddressService) Create(ctx context.Context, opt *UserAddressCreateOption) (address *models.UserAddress, err error) {

	//this.rep.Create(ctx)
	panic("")
}
