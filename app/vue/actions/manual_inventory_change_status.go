package actions

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/vue/actions"
)

type ManualInventoryChangeStatus struct {
	*actions.Action
}

func NewManualInventoryChangeStatus() *ManualInventoryChangeStatus {
	return &ManualInventoryChangeStatus{
		actions.NewAction(),
	}
}

func (this *ManualInventoryChangeStatus) Name() string {
	return "确认提交"
}

func (this *ManualInventoryChangeStatus) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (this *ManualInventoryChangeStatus) CanRun(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool {
	return true
}

func (this *ManualInventoryChangeStatus) ShowOnIndex() bool {
	return false
}

func (this *ManualInventoryChangeStatus) ShowOnDetail() bool {
	return true
}

func (this *ManualInventoryChangeStatus) HttpHandle(ctx *gin.Context, data map[string]interface{}, models []interface{}) error {
	spew.Dump(data)

	spew.Dump(models)

	return nil
}
