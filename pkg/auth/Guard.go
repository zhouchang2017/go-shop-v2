package auth

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Guard interface {
	// 确定当前用户是否经过身份验证
	Check() bool
	// 当前登录用户
	User() (user Authenticatable, err error)
	// 获取当前经过身份验证的用户的ID
	Id() (id string, err error)
	// 验证用户的凭证
	Validate(credentials map[string]string) bool
	// 设置上下文
	SetContext(ctx *gin.Context)
	// 获取上下文
	GetContext() *gin.Context
}

type RefreshToken interface {
	Refresh(duration time.Duration)
}
