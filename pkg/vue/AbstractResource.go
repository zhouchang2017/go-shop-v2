package vue

import (
	"github.com/gin-gonic/gin"
)

// 资源基本结构
type AbstractResource struct {
	Root *Vue
}

func (this *AbstractResource) SetRoot(vue *Vue) {
	this.Root = vue
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
func (this AbstractResource) Group() string {
	return "App"
}

// 是否显示在导航栏
func (this *AbstractResource) DisplayInNavigation(ctx *gin.Context) bool {
	return true
}

// 是否有列表页api
func (this *AbstractResource) ResourceHttpIndex() bool {
	return true
}

// 是否有详情页api
func (this *AbstractResource) ResourceHttpShow() bool {
	return true
}

// 是否有更新api
// 需要实现 ResourceHttpUpdate接口
func (this *AbstractResource) ResourceHttpUpdate() bool {
	return true
}

// 是否有创建api
// 需要实现 ResourceHttpCreate接口
func (this *AbstractResource) ResourceHttpCreate() bool {
	return true
}

// 是否有删除api
func (this *AbstractResource) ResourceHttpDelete() bool {
	return true
}

// 是否有硬删除api
func (this *AbstractResource) ResourceHttpForceDelete() bool {
	return true
}

// 是否有恢复api
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

// 自定义聚合
func (this *AbstractResource) Lenses() []Lens {
	return []Lens{}
}

// 自定义Link
func (this *AbstractResource) Links() []Link  {
	return []Link{}
}

// 自定义Cards
// 暂未写
