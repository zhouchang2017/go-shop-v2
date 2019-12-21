package ctx

import (
	"context"
	"github.com/gin-gonic/gin"
)

func GinCtxWithTrashed(ctx *gin.Context) {
	queryTrashed := ctx.Query("trashed")
	var trashed bool
	if queryTrashed == "true" {
		trashed = true
	}
	ctx.Set("trashed", trashed)
}

func GinCtxWithUser(ctx *gin.Context, user interface{}) *gin.Context {
	ctx.Set("user", user)
	return ctx
}

func GinCtxWithGuard(ctx *gin.Context, guard interface{}) *gin.Context {
	ctx.Set("guard", guard)
	return ctx
}

func WithTrashed(ctx context.Context, trashed bool) context.Context {
	return context.WithValue(ctx, "trashed", trashed)
}

func GetTrashed(ctx context.Context) bool {
	if value := ctx.Value("trashed"); value != nil {
		return value.(bool)
	}
	return false
}

func WithUser(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, "user", user)
}

func GetUser(ctx context.Context) interface{} {
	if value := ctx.Value("user"); value != nil {
		return value
	}
	return nil
}

func WithGuard(ctx context.Context, guard interface{}) context.Context {
	return context.WithValue(ctx, "guard", guard)
}

func GetGuard(ctx context.Context) interface{} {
	if value := ctx.Value("guard"); value != nil {
		return value
	}
	return nil
}

func WithForce(ctx context.Context, force bool) context.Context {
	return context.WithValue(ctx, "force", force)
}

func GetForce(ctx context.Context) bool {
	if value := ctx.Value("force"); value != nil {
		return value.(bool)
	}
	return false
}
