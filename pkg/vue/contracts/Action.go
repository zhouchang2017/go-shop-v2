package contracts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
)

// 动作
type Action interface {
	// 名称
	Name() string
	// 是否有权限可见
	AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool
	// 是否可执行
	CanRun(ctx *gin.Context,user auth.Authenticatable,model interface{}) bool
	// 列表页是否可见
	ShowOnIndex() bool
	// 详情页是否可见
	ShowOnDetail() bool
	// 列表页table row 可见，暂时不实现
	//ShowOnTableRow() bool
	// 处理函数
	HttpHandle(ctx *gin.Context, data map[string]interface{}, models []interface{}) error
}

// 执行动作提示文字
type ActionConfirmText interface {
	ConfirmText() string
}

// 动作表单字段
type ActionFields interface {
	Fields(ctx *gin.Context) []Field
}

// 动作确认按钮文字
type ActionConfirmButtonText interface {
	ConfirmButtonText() string
}

// 动作取消按钮文字
type ActionCancelButtonText interface {
	CancelButtonText() string
}
