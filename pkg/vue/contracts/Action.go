package contracts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/vue/message"
)

// 列表页动作操作：
// 导出excel、批量发送邮件、推送。。。
// 列表页动作应该和详情页动作区分

// 动作
type Action interface {
	Element
	// 名称
	Name() string
	// 是否可执行
	CanRun(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
	// 列表页是否可见
	ShowOnIndex() bool
	// 详情页是否可见
	ShowOnDetail() bool
	// 列表页table row 可见，暂时不实现
	//ShowOnTableRow() bool
	// 处理函数
	HttpHandle(ctx *gin.Context, data map[string]interface{}) (msg message.Message,err error)
	// 执行动作提示文字
	ConfirmText() string
	// 动作表单字段
	Fields(ctx *gin.Context) []Field
	// 动作确认按钮文字
	ConfirmButtonText() string
	// 动作取消按钮文字
	CancelButtonText() string
}
