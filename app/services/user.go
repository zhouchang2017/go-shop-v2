package services

import (
	"context"
	"github.com/medivhzhan/weapp/v2"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/qiniu"
	"time"
)

type UserService struct {
	rep *repositories.UserRep
}

func NewUserService(rep *repositories.UserRep) *UserService {
	return &UserService{rep: rep}
}

func (this *UserService) RegisterByWechat(ctx context.Context, info *weapp.UserInfo) (user *models.User, err error) {
	model := &models.User{
		WechatMiniId:  info.OpenID,
		WechatUnionId: info.UnionID,
		Nickname:      info.Nickname,
		Avatar:        qiniu.NewImage(info.Avatar),
		Gender:        info.Gender,
		Birth:         time.Time{},
		Country:       info.Country,
		Province:      info.Province,
		City:          info.City,
	}
	created := <-this.rep.Create(ctx, model)
	if created.Error != nil {
		return nil, created.Error
	}
	return created.Result.(*models.User),nil
}
