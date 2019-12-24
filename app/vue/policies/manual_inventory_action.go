package policies

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/auth"
)

func init() {
	register(NewManualInventoryActionPolicy)
}

type ManualInventoryActionPolicy struct {

}



func NewManualInventoryActionPolicy() *ManualInventoryActionPolicy {
	return &ManualInventoryActionPolicy{}
}

func (*ManualInventoryActionPolicy) Create(ctx *gin.Context, user auth.Authenticatable) bool {
	// 验证用户是否有门店，如果没有门店则不允许添加库存
	if admin, ok := user.(*models.Admin); ok {
		if len(admin.Shops) > 0 {
			return true
		}
	}
	return false
}