package contracts

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
)

type (
	// 资源接口
	Resource interface {
		// 资源名称
		Title() string
		// 侧边栏分组
		Group() string
		// 导航栏是否显示
		DisplayInNavigation(ctx *gin.Context, user interface{}) bool
		// 前端路由
		// 是否有列表页路由
		HasIndexRoute(ctx *gin.Context, user interface{}) bool
		// 是否有详情页路由
		HasDetailRoute(ctx *gin.Context, user interface{}) bool
		// 是否有更新页路由
		HasEditRoute(ctx *gin.Context, user interface{}) bool
		// 权限策略
		Policy() interface{}
		// 字段
		Fields(ctx *gin.Context, model interface{}) func() []interface{}
		// 实例对象
		Model() interface{}
		// 返回新资源
		Make(mode interface{}) Resource
		// 设置model
		SetModel(model interface{})
		// 聚合
		Lenses() []Lens
		// 自定义页面
		Pages() []Page
		// 过滤
		Filters(ctx *gin.Context) []Filter
		// 动作
		Actions(ctx *gin.Context) []Action
		// cards
		Cards(ctx *gin.Context) []Card
	}

	// 可展示icon图标
	Iconable interface {
		Icon() string
	}

	// 自定义uri
	CustomUri interface {
		UriKey() string
	}

	// 自定义列表页组件
	ResourceCustomIndexComponent interface {
		IndexComponent() Page
	}

	// 自定义详情页组件
	ResourceCustomDetailComponent interface {
		DetailComponent() Page
	}

	// 自定义创建页组件
	ResourceCustomCreationComponent interface {
		CreationComponent() Page
	}

	// 自定义更新页组件
	ResourceCustomUpdateComponent interface {
		UpdateComponent() Page
	}

	// 资源列表接口
	ResourcePaginationable interface {
		// 资源列表方法
		Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error)
	}

	// 资源详情接口
	ResourceShowable interface {
		// 资源详情页方法
		Show(ctx *gin.Context, id string) (res interface{}, err error)
	}

	// 资源创建接口
	ResourceStorable interface {
		// 资源创建方法
		Store(ctx *gin.Context, data map[string]interface{}) (redirect string, err error)
	}

	// 资源更新接口
	ResourceUpgradeable interface {
		// 资源更新方法
		Update(ctx *gin.Context, model interface{}, data map[string]interface{}) (redirect string, err error)
	}

	// 资源恢复接口
	ResourceRestoreable interface {
		// 资源恢复方法
		Restore(ctx *gin.Context, id string) (err error)
	}

	// 资源删除接口
	ResourceDestroyable interface {
		// 资源删除方法
		Destroy(ctx *gin.Context, id string) (err error)
	}

	// 资源硬删除接口
	ResourceForceDestroyable interface {
		// 资源硬删除方法
		ForceDestroy(ctx *gin.Context, id string) (err error)
	}
)
