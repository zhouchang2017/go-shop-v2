package vue

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/err"
	"reflect"
)

// 资源权限管理

// 列表api权限接口
type AuthorizedToViewAny interface {
	ViewAny(ctx *gin.Context, user auth.Authenticatable) bool
}

// 详情api权限接口
type AuthorizedToView interface {
	View(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

// 创建api权限接口
type AuthorizedToCreate interface {
	Create(ctx *gin.Context, user auth.Authenticatable) bool
}

// 更新api权限接口
type AuthorizedToUpdate interface {
	Update(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

// 删除api权限接口
type AuthorizedToDelete interface {
	Delete(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

// 硬删除api权限接口
type AuthorizedToForceDelete interface {
	ForceDelete(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

// 还原api权限接口
type AuthorizedToRestore interface {
	Restore(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

// 是否有创建权限
func (this *ResourceWarp) AuthorizedToCreate(ctx *gin.Context) bool {
	if creatable, ok := this.resource.(ResourceHttpCreate); ok && creatable.ResourceHttpCreate() {
		if policy, b := this.root.resolvePolicy(this.resource); b {
			if creatable, ok := policy.(AuthorizedToCreate); ok {
				return creatable.Create(ctx, ctx2.GetUser(ctx).(auth.Authenticatable))
			}
		}
		return true
	}
	return false
}

// 是否有删除权限
func (this *ResourceWarp) AuthorizedToDelete(ctx *gin.Context) (ok bool, model interface{}) {
	if r, ok := this.resource.(ResourceHttpDelete); ok && r.ResourceHttpDelete() {
		if policy, b := this.root.resolvePolicy(this.resource); b {
			if deletable, ok := policy.(AuthorizedToDelete); ok {
				result := <-this.resource.Repository().FindById(ctx, ctx.Param(this.SingularLabel()))
				if result.Error != nil {
					err.ErrorEncoder(nil, result.Error, ctx.Writer)
					ctx.Abort()
					return false, nil
				}
				this.resource.SetModel(result.Result)
				return deletable.Delete(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), result.Result), result.Result
			}
		}
		return true, nil
	}
	return false, nil
}

// 是否有硬删除权限
func (this *ResourceWarp) AuthorizedToForceDelete(ctx *gin.Context) (ok bool, model interface{}) {
	if r, ok := this.resource.(ResourceHttpForceDelete); ok && r.ResourceHttpForceDelete() {
		if policy, b := this.root.resolvePolicy(this.resource); b {
			if deletable, ok := policy.(AuthorizedToForceDelete); ok {
				result := <-this.resource.Repository().FindById(ctx, ctx.Param(this.SingularLabel()))
				if result.Error != nil {
					err.ErrorEncoder(nil, result.Error, ctx.Writer)
					ctx.Abort()
					return false, nil
				}
				this.resource.SetModel(result.Result)
				return deletable.ForceDelete(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), result.Result), result.Result
			}
		}
		return true, nil
	}
	return false, nil
}

// 是否有恢复权限
func (this *ResourceWarp) AuthorizedToRestore(ctx *gin.Context) (ok bool, model interface{}) {
	if r, ok := this.resource.(ResourceHttpRestore); ok && r.ResourceHttpRestore() {
		if policy, b := this.root.resolvePolicy(this.resource); b {
			if reformable, ok := policy.(AuthorizedToRestore); ok {
				result := <-this.resource.Repository().FindById(ctx, ctx.Param(this.SingularLabel()))
				if result.Error != nil {
					err.ErrorEncoder(nil, result.Error, ctx.Writer)
					ctx.Abort()
					return false, nil
				}
				this.resource.SetModel(result.Result)
				return reformable.Restore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), result.Result), result.Result
			}
		}
		return true, nil
	}
	return false, nil
}

// 是否有更新权限
func (this *ResourceWarp) AuthorizedToUpdate(ctx *gin.Context) (ok bool, model interface{}) {
	if upgradeable, ok := this.resource.(ResourceHttpUpdate); ok && upgradeable.ResourceHttpUpdate() {
		if policy, b := this.root.resolvePolicy(this.resource); b {
			if upgradeable, ok := policy.(AuthorizedToUpdate); ok {
				result := <-this.resource.Repository().FindById(ctx, ctx.Param(this.SingularLabel()))
				if result.Error != nil {
					err.ErrorEncoder(nil, result.Error, ctx.Writer)
					ctx.Abort()
					return false, nil
				}
				this.resource.SetModel(result.Result)
				return upgradeable.Update(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), result.Result), result.Result
			}
		}
		return true, nil
	}
	return false, nil
}

// 是否有详情权限
func (this *ResourceWarp) AuthorizedToView(ctx *gin.Context) (ok bool, model interface{}) {
	if r, ok := this.resource.(ResourceHttpShow); ok && r.ResourceHttpShow() {
		if policy, b := this.root.resolvePolicy(this.resource); b {
			if viewable, ok := policy.(AuthorizedToView); ok {
				result := <-this.resource.Repository().FindById(ctx, ctx.Param(this.SingularLabel()))
				if result.Error != nil {
					err.ErrorEncoder(nil, result.Error, ctx.Writer)
					ctx.Abort()
					return false, nil
				}
				return viewable.View(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), result.Result), result.Result
			}
		}
		return true, nil
	}
	return false, nil
}

// 是否有列表页权限
func (this *ResourceWarp) AuthorizedToViewAny(ctx *gin.Context) bool {
	if r, ok := this.resource.(ResourceHttpIndex); ok && r.ResourceHttpIndex() {
		if policy, b := this.root.resolvePolicy(this.resource); b {
			if viewAmiable, ok := policy.(AuthorizedToViewAny); ok {
				return viewAmiable.ViewAny(ctx, ctx2.GetUser(ctx).(auth.Authenticatable))
			}
		}
		return true
	}
	return false
}

func (this *ResourceWarp) AuthorizedTo(ctx *gin.Context, method string) bool {
	if policy, b := this.root.resolvePolicy(this.resource); b {
		of := reflect.ValueOf(policy)
		name := of.MethodByName(method)
		if !name.IsValid() {
			return false
		}
		call := name.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(ctx2.GetUser(ctx).(auth.Authenticatable))})
		if len(call) > 0 {
			return call[0].Interface().(bool)
		}
	}
	return true
}

func (this *ResourceWarp) Authorized(ctx *gin.Context, method string) {
	if !this.AuthorizedTo(ctx, method) {
		ctx.AbortWithStatus(403)
		return
	}
}
