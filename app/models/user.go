package models

import (
	"go-shop-v2/pkg/db/model"
	"go-shop-v2/pkg/qiniu"
	"time"
)

type User struct {
	model.MongoModel `inline`
	WechatMiniId     string      `json:"wechat_mini_id" bson:"wechat_mini_id"`
	WechatUnionId    string      `json:"wechat_union_id" bson:"wechat_union_id"`
	Nickname         string      `json:"nickname" bson:"nickname"`
	Avatar           qiniu.Image `json:"avatar" bson:"avatar"`
	Gender           int         `json:"gender" bson:"gender"`
	Birth            time.Time   `json:"birth" bson:"birth,omitempty"`
	Country          string      `json:"country" bson:"country"`
	Province         string      `json:"province" bson:"province"`
	City             string      `json:"city" bson:"city"`
}

type AssociatedUser struct {
	Id       string      `json:"id"`
	Nickname string      `json:"nickname"`
	Avatar   qiniu.Image `json:"avatar"`
	Gender   int         `json:"gender"`
}

func (user *User) ToAssociated() *AssociatedUser {
	return &AssociatedUser{
		Id:       user.GetID(),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
	}
}
