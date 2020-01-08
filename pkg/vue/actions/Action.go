package actions

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
	"go-shop-v2/pkg/vue/contracts"
)

type Action struct {
	meta            map[string]interface{}
	prefixComponent bool
}

func NewAction() *Action {
	return &Action{
		meta: map[string]interface{}{},
	}
}

func (Action) Component() string {
	return "action-dialog"
}

func (Action) PrefixComponent() bool {
	return false
}

func (this *Action) WithMeta(key string, value interface{}) {
	this.meta[key] = value
}

func (this Action) Meta() map[string]interface{} {
	return this.meta
}

func (Action) CancelButtonText() string {
	return "取消"
}

func (Action) ConfirmButtonText() string {
	return "确定"
}

func (Action) Fields(ctx *gin.Context) []contracts.Field {
	return []contracts.Field{}
}

func (Action) ConfirmText() string {
	return "是否确定执行该操作?"
}

func (Action) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	return true
}

func (Action) CanRun(ctx *gin.Context, user auth.Authenticatable, model interface{}) bool {
	return true
}

func (Action) ShowOnIndex() bool {
	return false
}

func (Action) ShowOnDetail() bool {
	return true
}
