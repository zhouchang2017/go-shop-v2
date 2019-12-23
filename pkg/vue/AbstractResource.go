package vue

import (
	"github.com/gin-gonic/gin"
)

// 资源基本结构
type AbstractResource struct {
}

// 是否有列表页路由
func (this *AbstractResource) HasIndexRoute(ctx *gin.Context) bool {
	return true
}

// 是否有详情页路由
func (this *AbstractResource) HasDetailRoute(ctx *gin.Context) bool {
	return true
}

// 是否有编辑页面路由
func (this *AbstractResource) HasEditRoute(ctx *gin.Context) bool {
	return true
}

// 左侧导航栏分组
func (this *AbstractResource) Group() string {
	return "App"
}

// 是否显示在导航栏
func (this *AbstractResource) DisplayInNavigation(ctx *gin.Context) bool {
	return true
}

func (this *AbstractResource) ResourceHttpIndex() bool {
	return true
}

func (this *AbstractResource) ResourceHttpShow() bool {
	return true
}

func (this *AbstractResource) ResourceHttpUpdate() bool {
	return true
}

func (this *AbstractResource) ResourceHttpCreate() bool {
	return true
}

func (this *AbstractResource) ResourceHttpDelete() bool {
	return true
}

func (this *AbstractResource) ResourceHttpForceDelete() bool {
	return true
}

func (this *AbstractResource) ResourceHttpRestore() bool {
	return true
}

// 创建成功钩子
func (this *AbstractResource) Created(ctx *gin.Context, resource interface{}) {

}

// 更新成功钩子
func (this *AbstractResource) Updated(ctx *gin.Context, resource interface{}) {

}

// 删除成功钩子
func (this *AbstractResource) Deleted(ctx *gin.Context, id string) {

}

// 恢复成功钩子
func (this *AbstractResource) Restored(ctx *gin.Context, resource interface{}) {

}

// 销毁成功钩子
func (this *AbstractResource) ForceDeleted(ctx *gin.Context, id string) {

}
