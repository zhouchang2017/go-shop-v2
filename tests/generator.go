package tests

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/model"
	"time"
)

func GenerateUser() *models.User {
	mongoSrt := model.MongoModel{}
	mongoSrt.SetID("testuserid")
	//user
	authUser := &models.User{
		MongoModel:    mongoSrt,
		WechatMiniId:  "test_user_wechat_mini_id",
		WechatUnionId: "test_user_wechat_union_id",
		Nickname:      "test_user_nickname",
		Avatar:        "test_user_avatar",
		Gender:        1,
		Birth:         time.Now(),
		Country:       "test_user_country",
		Province:      "test_user_province",
		City:          "test_user_city",
	}
	// return
	return authUser
}
