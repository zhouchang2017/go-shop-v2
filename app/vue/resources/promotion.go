package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/app/vue/pages"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

type Promotion struct {
	core.AbstractResource
	model   interface{}
	service *services.PromotionService
}

func (this *Promotion) Destroy(ctx *gin.Context, id string) (err error) {
	return this.service.Delete(ctx,id)
}

func NewPromotionResource() *Promotion {
	return &Promotion{model: &models.Promotion{}, service: services.MakePromotionService()}
}

// 自定义创建页
func (this *Promotion) CreationComponent() contracts.Page {
	return pages.NewPromotionCreatePage()
}

// 自定义更新页
func (this *Promotion) UpdateComponent() contracts.Page {
	return pages.NewPromotionUpdatePage()
}


// 实现列表页api
func (this *Promotion) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	req.SetSearchField("name")
	return this.service.Pagination(ctx, req)
}

// 实现详情页api
func (this *Promotion) Show(ctx *gin.Context, id string) (res interface{}, err error) {
	return this.service.FindById(ctx, id)
}

func (p Promotion) Title() string {
	return "促销管理"
}

func (p Promotion) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {
		return []interface{}{
			fields.NewIDField(),
			fields.NewTextField("活动名称", "Name"),
			fields.NewStatusField("类型", "Type").WithOptions([]*fields.StatusOption{
				fields.NewStatusOption("单品活动", 0),
				fields.NewStatusOption("复合活动", 1),
			}),
			fields.NewStatusField("互斥", "Mutex").WithOptions([]*fields.StatusOption{
				fields.NewStatusOption("是", true).Success(),
				fields.NewStatusOption("否", false).Cancel(),
			}),
			fields.NewStatusField("启用", "Enable").WithOptions([]*fields.StatusOption{
				fields.NewStatusOption("是", true).Success(),
				fields.NewStatusOption("否", false).Cancel(),
			}),
			fields.NewDateTime("开始时间", "BeginAt"),
			fields.NewDateTime("结束时间", "EndedAt"),

			panels.NewPanel("促销规则",
				fields.NewStatusField("类型", "Rule.Type", fields.OnlyOnDetail()).WithOptions([]*fields.StatusOption{
					fields.NewStatusOption("不限", 0),
					fields.NewStatusOption("金额大于", 1),
					fields.NewStatusOption("数量大于", 2),
				}, ),
				fields.NewTextField("值", "Rule.Value", fields.OnlyOnDetail()),
			),

			panels.NewPanel("优惠策略",
				fields.NewStatusField("类型", "Policy.Type", fields.OnlyOnDetail()).WithOptions([]*fields.StatusOption{
					fields.NewStatusOption("打折", 1),
					fields.NewStatusOption("直减", 2),
					fields.NewStatusOption("免邮", 3),
				}, ),
				fields.NewTextField("值", "Policy.Value", fields.OnlyOnDetail()),
			),

			fields.NewHasManyField("促销商品", &PromotionItemResource{}),
		}
	}
}

func (p *Promotion) Model() interface{} {
	return p.model
}

func (p Promotion) Make(mode interface{}) contracts.Resource {
	return &Promotion{
		model:   mode,
		service: p.service,
	}
}

func (p *Promotion) SetModel(model interface{}) {
	p.model = model
}

func (this Promotion) Group() string {
	return "Product"
}

func (this Promotion) Icon() string {
	return "icons-announcement"
}
