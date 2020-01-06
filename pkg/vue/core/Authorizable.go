package core

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/vue/contracts"
)

var AuthorizedBefore func(ctx *gin.Context, user auth.Authenticatable) bool

// 是否有创建权限
func AuthorizedToCreate(ctx *gin.Context, resource contracts.Resource) bool {

	// 是否实现自定义创建页
	if customCreation, ok := resource.(contracts.ResourceCustomCreationComponent); ok {

		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}

		if customCreation.CreationComponent().AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			return true
		}
	}
	// 是否实现创建方法
	if _, ok := resource.(contracts.ResourceStorable); ok {
		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}
		// 是否实现创建权限
		if policy, ok := resource.Policy().(contracts.AuthorizedToCreate); ok {
			return policy.Create(ctx, ctx2.GetUser(ctx).(auth.Authenticatable))
		}
		return true
	}
	return false
}

// 是否有删除权限
func AuthorizedToDelete(ctx *gin.Context, resource contracts.Resource) (ok bool) {

	// 是否实现删除方法
	if _, ok := resource.(contracts.ResourceDestroyable); ok {
		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}

		// 是否实现详情方法
		if showable, implement := resource.(contracts.ResourceShowable); implement {
			if policy, ok := resource.Policy().(contracts.AuthorizedToDelete); ok {
				model := resource.Model()
				if model == nil {
					res, err := showable.Show(ctx, ResourceIdParam(resource))
					if err != nil {
						err2.ErrorEncoder(ctx, err, ctx.Writer)
						ctx.Abort()
						return false
					}
					resource.SetModel(res)
					model = res
				}

				return policy.Delete(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), model)
			}
			return true
		}
	}
	return false
}

// 是否有硬删除权限
func AuthorizedToForceDelete(ctx *gin.Context, resource contracts.Resource) (ok bool) {

	// 是否实现删除方法
	if _, ok := resource.(contracts.ResourceForceDestroyable); ok {
		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}

		// 是否实现详情方法
		if showable, implement := resource.(contracts.ResourceShowable); implement {
			if policy, ok := resource.Policy().(contracts.AuthorizedToForceDelete); ok {
				model := resource.Model()
				if model == nil {
					res, err := showable.Show(ctx, ResourceIdParam(resource))
					if err != nil {
						err2.ErrorEncoder(ctx, err, ctx.Writer)
						ctx.Abort()
						return false
					}
					resource.SetModel(res)
					model = res
				}

				return policy.ForceDelete(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), model)
			}
			return true
		}
	}
	return false
}

// 是否有恢复权限
func AuthorizedToRestore(ctx *gin.Context, resource contracts.Resource) (ok bool) {

	// 是否实现恢复方法
	if _, ok := resource.(contracts.ResourceRestoreable); ok {
		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}

		// 是否实现详情方法
		if showable, implement := resource.(contracts.ResourceShowable); implement {
			if policy, ok := resource.Policy().(contracts.AuthorizedToRestore); ok {
				model := resource.Model()
				if model == nil {
					res, err := showable.Show(ctx, ResourceIdParam(resource))
					if err != nil {
						err2.ErrorEncoder(ctx, err, ctx.Writer)
						ctx.Abort()
						return false
					}
					resource.SetModel(res)
					model = res
				}
				return policy.Restore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), model)
			}
			return true
		}
	}
	return false
}

// 是否有更新权限
func AuthorizedToUpdate(ctx *gin.Context, resource contracts.Resource) (ok bool) {

	// 是否实现自定义更新页
	if customUpdate, ok := resource.(contracts.ResourceCustomUpdateComponent); ok {

		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}

		if customUpdate.UpdateComponent().AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			return true
		}
	}

	// 是否实现更新方法
	if _, ok := resource.(contracts.ResourceUpgradeable); ok {
		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}

		// 是否实现详情方法
		if showable, implement := resource.(contracts.ResourceShowable); implement {
			if policy, ok := resource.Policy().(contracts.AuthorizedToUpdate); ok {
				model := resource.Model()
				if model == nil {
					res, err := showable.Show(ctx, ResourceIdParam(resource))
					if err != nil {
						err2.ErrorEncoder(ctx, err, ctx.Writer)
						ctx.Abort()
						return false
					}
					resource.SetModel(res)
					model = res
				}
				return policy.Update(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), model)
			}
			return true
		}
	}
	return false
}

// 是否有详情权限
func AuthorizedToView(ctx *gin.Context, resource contracts.Resource) (ok bool) {

	// 是否实现自定义详情页
	if customDetail, ok := resource.(contracts.ResourceCustomDetailComponent); ok {

		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}

		if customDetail.DetailComponent().AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			return true
		}
	}

	// 是否实现详情方法
	if showable, implement := resource.(contracts.ResourceShowable); implement {
		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}

		if policy, ok := resource.Policy().(contracts.AuthorizedToView); ok {
			model := resource.Model()
			if model == nil {
				res, err := showable.Show(ctx, ResourceIdParam(resource))
				if err != nil {
					err2.ErrorEncoder(ctx, err, ctx.Writer)
					ctx.Abort()
					return false
				}
				resource.SetModel(res)
				model = res
			}
			return policy.View(ctx, ctx2.GetUser(ctx).(auth.Authenticatable), model)
		}
		return true
	}
	return false
}

// 是否有列表页权限
func AuthorizedToViewAny(ctx *gin.Context, resource contracts.Resource) bool {

	// 是否实现自定义列表页
	if customIndex, ok := resource.(contracts.ResourceCustomIndexComponent); ok {

		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}

		if customIndex.IndexComponent().AuthorizedTo(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
			return true
		}
	}

	// 是否实现详情方法
	if _, implement := resource.(contracts.ResourcePaginationable); implement {
		if AuthorizedBefore != nil {
			if AuthorizedBefore(ctx, ctx2.GetUser(ctx).(auth.Authenticatable)) {
				return true
			}
		}
		if policy, ok := resource.Policy().(contracts.AuthorizedToViewAny); ok {
			return policy.ViewAny(ctx, ctx2.GetUser(ctx).(auth.Authenticatable))
		}
		return true
	}
	return false
}
