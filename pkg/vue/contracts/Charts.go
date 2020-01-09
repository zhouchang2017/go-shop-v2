package contracts

import (
	"github.com/gin-gonic/gin"
)

// 图表
type Charts interface {
	Card
	Name() string      // 名称
	Columns() []string // 列
	Settings() map[string]interface{}
	Extend() map[string]interface{}
	// 处理函数
	// resourceName
	// resourceId
	HttpHandle(ctx *gin.Context) (rows interface{}, err error)
}
