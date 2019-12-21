package vue

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	"go-shop-v2/pkg/err"
)

// 资源权限管理

type AuthorizedToViewAny interface {
	ViewAny(ctx *gin.Context, user auth.Authenticatable) bool
}

type AuthorizedToView interface {
	View(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

type AuthorizedToCreate interface {
	Create(ctx *gin.Context, user auth.Authenticatable) bool
}

type AuthorizedToUpdate interface {
	Update(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

type AuthorizedToDelete interface {
	Delete(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

type AuthorizedToForceDelete interface {
	ForceDelete(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

type AuthorizedToRestore interface {
	Restore(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool
}

func (this *ResourceWarp) AuthorizedToCreate(ctx *gin.Context) bool {
	if policy, b := this.root.resolvePolicy(this.resource); b {
		if creatable, ok := policy.(AuthorizedToCreate); ok {
			return creatable.Create(ctx, ctx2.GetUser(ctx).(auth.Authenticatable))
		}
	}
	return true
}

func (this *ResourceWarp) AuthorizedToDelete(ctx *gin.Context) (ok bool, model interface{}) {
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

func (this *ResourceWarp) AuthorizedToForceDelete(ctx *gin.Context) (ok bool, model interface{}) {
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

func (this *ResourceWarp) AuthorizedToRestore(ctx *gin.Context) (ok bool, model interface{}) {
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

func (this *ResourceWarp) AuthorizedToUpdate(ctx *gin.Context) (ok bool, model interface{}) {
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

func (this *ResourceWarp) AuthorizedToView(ctx *gin.Context) (ok bool, model interface{}) {
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

func (this *ResourceWarp) AuthorizedToViewAny(ctx *gin.Context) bool {
	if policy, b := this.root.resolvePolicy(this.resource); b {
		if viewAmiable, ok := policy.(AuthorizedToViewAny); ok {
			return viewAmiable.ViewAny(ctx, ctx2.GetUser(ctx).(auth.Authenticatable))
		}
	}
	return true
}
