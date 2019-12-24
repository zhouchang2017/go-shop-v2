package vue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
)

// Lens
// 自定义查询聚合

// routerName lenses.resource.urikey
type Lens interface {
	Title() string     // 标题
	Component() string // 对应vue组件
	RouterName() string
	HttpHandle() func(ctx *gin.Context)                            // http处理
	AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool // 授权
}

func lensRouterName(lens Lens, resourceName string) string {
	return fmt.Sprintf("lenses.%s.%s", resourceName, lens.RouterName())
}

func lensApiUri(lens Lens, resourceUriKey string) string {
	return fmt.Sprintf("/lenses/%s/%s", resourceUriKey, lens.RouterName())
}
