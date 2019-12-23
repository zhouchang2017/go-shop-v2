package resources

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/app/services"
	ctx2 "go-shop-v2/pkg/ctx"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/request"
	"go-shop-v2/pkg/vue"
)

func init() {
	register(NewManualInventoryActionResource)
}

type ManualInventoryAction struct {
	vue.AbstractResource
	model   *models.ManualInventoryAction
	rep     *repositories.ManualInventoryActionRep
	service *services.ManualInventoryActionService
}

func NewManualInventoryActionResource(rep *repositories.ManualInventoryActionRep, service *services.ManualInventoryActionService) *ManualInventoryAction {
	return &ManualInventoryAction{model: &models.ManualInventoryAction{}, rep: rep, service: service}
}

type manualInventoryActionForm struct {
	Type   int8                                `json:"type" form:"type"`
	ShopId string                              `json:"shop_id" form:"shop_id"`
	Items  []*models.ManualInventoryActionItem `json:"items" form:"items"`
	Status int8                                `json:"status" form:"status"`
}

func (m *ManualInventoryAction) UpdateFormParse(ctx *gin.Context, model interface{}) (entity interface{}, err error) {
	form := &manualInventoryActionForm{}
	if err := ctx.ShouldBind(form); err != nil {
		return nil, err
	}
	action := model.(*models.ManualInventoryAction)
	user := ctx2.GetUser(ctx)
	if admin, ok := user.(*models.Admin); ok {
		if err := action.SetType(form.Type); err != nil {
			return nil, err
		}
		// 创建设置为保存状态
		action.SetStatusToSaved()
		action.Items = form.Items
		action.User = admin.ToAssociated()
		if form.ShopId != action.Shop.Id {
			return m.service.SetShop(ctx, action, form.ShopId)
		}
		return action, nil
	}
	return nil, err2.Err401
}

func (m *ManualInventoryAction) CreateFormParse(ctx *gin.Context) (entity interface{}, err error) {
	form := &manualInventoryActionForm{}
	if err := ctx.ShouldBind(form); err != nil {
		return nil, err
	}
	action := &models.ManualInventoryAction{}
	user := ctx2.GetUser(ctx)
	if admin, ok := user.(*models.Admin); ok {
		if err := action.SetType(form.Type); err != nil {
			return nil, err
		}
		// 创建设置为保存状态
		action.SetStatusToSaved()
		action.Items = form.Items
		action.User = admin.ToAssociated()
		return m.service.SetShop(ctx, action, form.ShopId)
	}
	return nil, err2.Err401
}

func (m *ManualInventoryAction) IndexQuery(ctx *gin.Context, request *request.IndexRequest) error {
	return nil
}

func (m *ManualInventoryAction) Model() interface{} {
	return m.model
}

func (m *ManualInventoryAction) Repository() repository.IRepository {
	return m.rep
}

func (m ManualInventoryAction) Make(model interface{}) vue.Resource {
	return &ManualInventoryAction{model: model.(*models.ManualInventoryAction)}
}

func (m *ManualInventoryAction) SetModel(model interface{}) {
	m.model = model.(*models.ManualInventoryAction)
}

func (m ManualInventoryAction) Title() string {
	return "库存操作"
}

func (ManualInventoryAction) Group() string {
	return "Shop"
}

func (ManualInventoryAction) Icon() string {
	return "i-repeat"
}