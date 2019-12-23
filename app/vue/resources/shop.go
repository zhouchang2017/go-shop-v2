package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/events"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/event"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue"
)

func init() {
	register(NewShopResource)
}

type Shop struct {
	vue.AbstractResource
	model   *models.Shop
	rep     *repositories.ShopRep
	service *services.ShopService
}

func NewShopResource(rep *repositories.ShopRep, service *services.ShopService) *Shop {
	return &Shop{model: &models.Shop{}, rep: rep, service: service}
}

type shopForm struct {
	Name     string              `json:"name" form:"name" binding:"required"`
	Address  *models.ShopAddress `json:"address"`           // 地址
	Location *models.Location    `json:"location"`          // 坐标
	Members  []string            `json:"members,omitempty"` // 成员
}

// 更新表单处理
func (s *Shop) UpdateFormParse(ctx *gin.Context, model interface{}) (entity interface{}, err error) {
	shop := model.(*models.Shop)
	form := &shopForm{}
	if err = ctx.ShouldBind(form); err != nil {
		return nil, err
	}
	shop.Name = form.Name
	shop.Location = form.Location
	shop.Address = form.Address

	return s.service.SetMembers(ctx, shop, form.Members...)
}

// 创建表单处理
func (s *Shop) CreateFormParse(ctx *gin.Context) (entity interface{}, err error) {
	form := &shopForm{}
	if err = ctx.ShouldBind(form); err != nil {
		return nil, err
	}
	shop := &models.Shop{
		Name:     form.Name,
		Address:  form.Address,
		Location: form.Location,
	}

	return s.service.SetMembers(ctx, shop, form.Members...)
}

// 创建成功钩子
func (s *Shop) Created(ctx *gin.Context, resource interface{}) {
	event.Dispatch(events.ShopCreated{Shop: resource.(*models.Shop)})
}

// 更新成功钩子
func (s *Shop) Updated(ctx *gin.Context, resource interface{}) {
	event.Dispatch(events.ShopUpdated{Shop: resource.(*models.Shop)})
}

// 列表页搜索处理
func (s *Shop) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	return nil
}

func (s *Shop) Model() interface{} {
	return s.model
}

func (s *Shop) Repository() repository.IRepository {
	return s.rep
}

func (s Shop) Make(model interface{}) vue.Resource {
	return &Shop{model: model.(*models.Shop)}
}

func (s *Shop) SetModel(model interface{}) {
	s.model = model.(*models.Shop)
}

func (s Shop) Title() string {
	return "门店"
}

func (Shop) Icon() string {
	return "i-store"
}

func (Shop) Group() string {
	return "Shop"
}
