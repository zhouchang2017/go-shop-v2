package wechat

import (
	"github.com/medivhzhan/weapp/v2"
	"time"
)

// getAllDelivery
// 获取物流公司
func (this *sdk) GetAllDelivery() (*weapp.DeliveryList, error) {
	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}
	delivery, err := weapp.GetAllDelivery(token)
	if err != nil {
		return nil, err
	}
	if err := delivery.GetResponseError(); err != nil {
		return nil, err
	}
	return delivery, nil
}

// getAllAccount
// 获取所有绑定的物流账号
func (this *sdk) GetAllAccount() (*weapp.AccountList, error) {
	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}
	account, err := weapp.GetAllAccount(token)
	if err != nil {
		return nil, err
	}
	if err := account.GetResponseError(); err != nil {
		// 微信方异常处理
		return nil, err
	}
	return account, nil
}

type CreateExpressOrderOption struct {
	*weapp.ExpressOrder
	AddSource  weapp.ExpressOrderSource
	ExpectTime uint
}

// addOrder
// 生成运单
func (this *sdk) AddExpressOrder(opt *CreateExpressOrderOption) (response *weapp.CreateExpressOrderResponse, err error) {
	if !this.config.IsProd {
		opt.DeliveryID = "TEST"
		opt.BizID = "test_biz_id"
		opt.Service.Type = 1
		opt.Service.Name = "test_service_name"
	}
	creator := weapp.ExpressOrderCreator{
		ExpressOrder: *opt.ExpressOrder,
		AddSource:    opt.AddSource,
		WXAppID:      this.config.AppId,
		ExpectTime:   opt.ExpectTime,
	}
	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}
	if response, err = creator.Create(token); err != nil {
		return nil, err
	}

	if err = response.GetResponseError(); err != nil {
		return nil, err
	}
	return
}

type CancelExpressOrderOption struct {
	OrderId    string `json:"order_id"`
	OpenId     string `json:"open_id"`
	DeliveryId string `json:"delivery_id"`
	WaybillId  string `json:"waybill_id"`
}

// cancelOrder
// 取消运单
func (this *sdk) CancelOrder(opt *CancelExpressOrderOption) (*weapp.CommonError, error) {
	canceler := weapp.ExpressOrderCanceler{
		OrderID:    opt.OrderId,
		OpenID:     opt.OpenId,
		DeliveryID: opt.DeliveryId,
		WaybillID:  opt.WaybillId,
	}
	if !this.config.IsProd {
		canceler.DeliveryID = "TEST"
	}

	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}

	return canceler.Cancel(token)
}

type GetterExpressOrderOption struct {
	OrderID    string
	OpenID     string
	DeliveryID string
	WaybillID  string
}

// getOrder
// 获取运单数据
func (this *sdk) GetOrder(opt GetterExpressOrderOption) (response *weapp.GetExpressOrderResponse, err error) {
	getter := weapp.ExpressOrderGetter{
		OrderID:    opt.OrderID,
		OpenID:     opt.OpenID,
		DeliveryID: opt.DeliveryID,
		WaybillID:  opt.WaybillID,
	}
	if !this.config.IsProd {
		getter.DeliveryID = "TEST"
	}

	var token string

	token, err = this.getAccessToken()
	if err != nil {
		return nil, err
	}

	response, err = getter.Get(token)
	if err != nil {
		return nil, err
	}
	if err := response.GetResponseError(); err != nil {
		return nil, err
	}
	return
}

type TestExpressUpdateOrderOption struct {
	OrderId    string
	WaybillID  string
	ActionTime time.Time
	ActionType uint
	ActionMsg  string
}

// testUpdateOrder
// 模拟快递公司更新订单状态, 该接口只能用户测试
func (this *sdk) TestUpdateOrder(opt *TestExpressUpdateOrderOption) (*weapp.CommonError, error) {

	updateExpressOrderTester := weapp.UpdateExpressOrderTester{
		BizID:      "test_biz_id",
		OrderID:    opt.OrderId,
		WaybillID:  opt.WaybillID,
		DeliveryID: "TEST",
		ActionTime: uint(opt.ActionTime.Unix()),
		ActionType: opt.ActionType,
		ActionMsg:  opt.ActionMsg,
	}

	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}

	return updateExpressOrderTester.Test(token)
}

// onPathUpdate
// 运单轨迹更新事件。当运单轨迹有更新时，会产生如下数据包。收到事件之后，回复success或者空串即可
func (this *sdk) OnPathUpdate() {

}
