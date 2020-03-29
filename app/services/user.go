package services

import (
	"context"
	"github.com/medivhzhan/weapp/v2"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type UserService struct {
	rep *repositories.UserRep
}

func (this *UserService) RetrieveById(identifier interface{}) (auth.Authenticatable, error) {
	return this.rep.RetrieveById(identifier)
}

func (this *UserService) RetrieveByCredentials(credentials map[string]string) (auth.Authenticatable, error) {
	return this.rep.RetrieveByCredentials(credentials)
}

func (this *UserService) ValidateCredentials(user auth.Authenticatable, credentials map[string]string) bool {
	return this.rep.ValidateCredentials(user, credentials)
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
	return created.Result.(*models.User), nil
}

// 获取当天新注册用户数
func (this *UserService) TodayNewUserCount(ctx context.Context) (count int64) {
	result := <-this.rep.Count(ctx,
		bson.M{
			"created_at": bson.M{"$gte": utils.TodayStart(), "$lte": utils.TodayEnd()},
		})
	if result.Error != nil {
		return 0
	}
	return result.Result
}
