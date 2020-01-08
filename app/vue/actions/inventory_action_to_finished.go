package actions

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/vue/actions"
	"go-shop-v2/pkg/vue/message"
)

// 库存操作，状态改变为finished动作
type InventoryActionToFinished struct {
	*actions.Action
	service *services.ManualInventoryActionService
}

func NewInventoryActionToFinished() *InventoryActionToFinished {
	return &InventoryActionToFinished{
		Action:  actions.NewAction(),
		service: services.MakeManualInventoryActionService(),
	}
}

// 动作名称
func (this *InventoryActionToFinished) Name() string {
	return "确认提交"
}

// 权限
func (this *InventoryActionToFinished) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

// 是否可以运行
func (this *InventoryActionToFinished) CanRun(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool {
	inventoryAction := model.(*models.ManualInventoryAction)
	// 只有为未提交状态才能运行
	return inventoryAction.IsSaved()
}

// 列表页不显示，（列表页动作暂未实现）
func (this *InventoryActionToFinished) ShowOnIndex() bool {
	return false
}

// 详情页显示
func (this *InventoryActionToFinished) ShowOnDetail() bool {
	return true
}

type formData struct {
	Resources []string `json:"resources"`
}

// 处理函数
func (this *InventoryActionToFinished) HttpHandle(ctx *gin.Context, data map[string]interface{}) (msg message.Message, err error) {
	form := formData{}
	err = mapstructure.Decode(data, &form)
	if err != nil {
		return
	}

	for _, id := range form.Resources {
		_, err = this.service.StatusToFinished(ctx, id)
		if err != nil {
			return
		}
	}
	return message.Success("提交成功"), nil
}

// 动作可以包含自定义字段