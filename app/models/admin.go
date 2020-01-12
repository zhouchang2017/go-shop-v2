package models

import (
	"fmt"
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/utils"
)

var adminTypes = []string{"root", "admin", "manager", "salesman"}


// 后台用户
// 用户类型：root、admin、manager、salesman
type Admin struct {
	model.MongoModel `inline`
	Username         string            `json:"username"`
	Password         string            `json:"-"`
	Nickname         string            `json:"nickname"`
	Type             string            `json:"type"`
	Shops            []*AssociatedShop `json:"shops" bson:"shops"`
}

// 关联简单管理员结构
type AssociatedAdmin struct {
	Id       string `json:"id"`
	Nickname string `json:"nickname"`
	Type     string `json:"type"`
}

func NewAdmin() *Admin {
	return &Admin{}
}

func (a Admin) ToAssociated() *AssociatedAdmin {
	return &AssociatedAdmin{
		Id:       a.GetID(),
		Nickname: a.Nickname,
		Type:     a.Type,
	}
}

func (a *Admin) SetType(t string) (*Admin, error) {
	if exist, _ := utils.InArray(t, adminTypes); exist {
		a.Type = t
		return a, nil
	}
	return a, fmt.Errorf("Type [%s] not allow!", t)
}

func (a *Admin) GetAuthIdentifierName() string {
	return "id"
}

func (a *Admin) GetAuthIdentifier() string {
	return a.GetID()
}

func (a *Admin) GetAuthPassword() string {
	return a.Password
}

func (a *Admin) GetNickname() string {
	return a.Nickname
}

func (a *Admin) IsManager() bool {
	return a.Type == "manager"
}

func (a *Admin) RandNickname() {
	a.Nickname = utils.RandomString(10)
}

func (a *Admin) SetPassword(password string) (err error) {
	a.Password, err = utils.Encrypt(password)
	return err
}

func (a *Admin) SetHashPassword(password string) {
	a.Password = password
}

//func (a *Admin) SetShopIds(shopIds []string) *Admin {
//	a.ShopIds = a.ShopIds.Make(shopIds)
//	return a
//}
//
//func (a *Admin) GetShopIds() []string {
//	return a.ShopIds.Split()
//}

func (a *Admin) SetShops(shops []*Shop) *Admin {
	a.Shops = []*AssociatedShop{}
	for _, shop := range shops {
		a.Shops = append(a.Shops, shop.ToAssociated())
	}
	return a
}

//func (this *Admin) LoadAddress() error {
//	this.Address = &Address{}
//	return this.BelongsTo(this, this.Address)
//}
