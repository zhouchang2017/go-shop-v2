package http

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"net/http"
)

// 用户地址
type AddressController struct {
	addressSrv *services.AddressService
}

// 用户地址列表（不做分页）
func (this *AddressController) Index(ctx *gin.Context) {
	user := ctx2.GetUser(ctx).(*models.User)

	addresses, err := this.addressSrv.Index(ctx, user.GetID())

	if err != nil {
		ResponseError(ctx, err)
		return
	}

	Response(ctx, addresses, http.StatusOK)
}

// 添加地址
func (this *AddressController) Add(ctx *gin.Context) {
	user := ctx2.GetUser(ctx).(*models.User)
	var form services.UserAddressCreateOption
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}
	form.UserId = user.GetID()
	address, err := this.addressSrv.Create(ctx, &form)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, address, http.StatusCreated)
}

// 更新地址
func (this *AddressController) Update(ctx *gin.Context) {
	user := ctx2.GetUser(ctx).(*models.User)
	id := ctx.Param("id")
	var form services.UserAddressCreateOption
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
		return
	}
	address, err := this.addressSrv.FindById(ctx, id)
	if err != nil {
		// 404 地址不存在
		ResponseError(ctx, err)
	}
	if form.Province != nil {
		address.Province = *form.Province
	}
	if form.City != nil {
		address.City = *form.City
	}
	if form.Areas != nil {
		address.Areas = *form.Areas
	}
	if form.Addr != nil {
		address.Addr = *form.Addr
	}
	if form.ContactName != nil {
		address.ContactName = *form.ContactName
	}
	if form.ContactPhone != nil {
		address.ContactPhone = *form.ContactPhone
	}
	address.IsDefault = form.IsDefault

	address.UserId = user.GetID()

	updated, err := this.addressSrv.Update(ctx, address)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, updated, http.StatusOK)
}

// 删除地址
func (this *AddressController) Delete(ctx *gin.Context) {
	user := ctx2.GetUser(ctx).(*models.User)
	id := ctx.Param("id")
	address, err := this.addressSrv.FindById(ctx, id)
	if err != nil {
		// 404
		ResponseError(ctx, err)
		return
	}
	if address.UserId != user.GetID() {
		// 异常
		ResponseError(ctx, err2.New(http.StatusForbidden, "无权删除"))
		return
	}
	if err := this.addressSrv.Delete(ctx, id); err != nil {
		ResponseError(ctx, err)
		return
	}
	Response(ctx, nil, http.StatusOK)
}
