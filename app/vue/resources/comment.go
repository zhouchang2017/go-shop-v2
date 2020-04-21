package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/response"
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/core"
	"go-shop-v2/pkg/vue/fields"
	"go-shop-v2/pkg/vue/panels"
)

type Comment struct {
	core.AbstractResource
	model   interface{}
	service *services.CommentService
}

func (c Comment) SearchPlaceholder() string {
	return "请输入订单号"
}

func NewCommentResource() *Comment {
	return &Comment{
		AbstractResource: core.AbstractResource{},
		model:            &models.Comment{},
		service:          services.MakeCommentService(),
	}
}

func (c Comment) Title() string {
	return "评论"
}

func (this *Comment) Pagination(ctx *gin.Context, req *request.IndexRequest) (res interface{}, pagination response.Pagination, err error) {
	req.SetSearchField("order_no")
	orderNo := ctx.Query("order_no")
	if orderNo != "" {
		req.AppendFilter("order_no", orderNo)
	}
	productId := ctx.Query("product_id")
	if productId != "" {
		req.AppendFilter("product_id", productId)

	}
	return this.service.Pagination(ctx, req)
}

func (c Comment) Fields(ctx *gin.Context, model interface{}) func() []interface{} {
	return func() []interface{} {

		return []interface{}{
			fields.NewIDField(fields.WithMeta("min-width", 150)),
			fields.NewTextField("商品ID", "ProductId", fields.WithMeta("min-width", 150)).Link(&Product{}, "ProductId"),
			fields.NewTextField("订单号", "OrderNo", fields.WithMeta("min-width", 150)),


			panels.NewPanel("用户",
				fields.NewAvatar("头像", "User.Avatar", fields.SetShowOnIndex(false)).RoundedFull(),
				fields.NewTextField("用户", "User.Nickname"),
			),

			fields.NewTextField("评分", "Rate"),
			fields.NewTextField("内容", "Content"),
			fields.NewDateTime("创建时间", "CreatedAt", fields.SetShowOnIndex(false)),
			fields.NewDateTime("更新时间", "UpdatedAt", fields.WithMeta("min-width", 150)),
		}
	}
}

func (c Comment) Model() interface{} {
	return c.model
}

func (c *Comment) Make(mode interface{}) contracts.Resource {
	return &Comment{
		AbstractResource: core.AbstractResource{},
		model:            mode,
		service:          c.service,
	}
}

func (c *Comment) SetModel(model interface{}) {
	c.model = model
}

func (r Comment) Icon() string {
	return "icons-comment"
}

func (r Comment) Group() string {
	return "Order"
}
