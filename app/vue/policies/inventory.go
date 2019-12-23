package policies

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/auth"
)

func init() {
	register(NewInventoryPolicy)
}

type InventoryPolicy struct {

}

func NewInventoryPolicy() *InventoryPolicy {
	return &InventoryPolicy{}
}

func (*InventoryPolicy) Create(ctx *gin.Context, user auth.Authenticatable) bool {
	spew.Dump(user)
	// 验证用户是否有门店，如果没有门店则不允许添加库存
	if admin, ok := user.(*models.Admin); ok {
		if len(admin.Shops) > 0 {
			return true
		}
	}
	return false
}
