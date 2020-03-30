package wechat

import (
	"github.com/medivhzhan/weapp/v2"
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
func (this *sdk) AddExpressOrder(opt CreateExpressOrderOption, sanbox bool) (response *weapp.CreateExpressOrderResponse, err error) {
	if sanbox {
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
// 取消运动
func (this *sdk) CancelOrder(opt CancelExpressOrderOption) (*weapp.CommonError, error) {
	canceler := weapp.ExpressOrderCanceler{
		OrderID:    opt.OrderId,
		OpenID:     opt.OpenId,
		DeliveryID: opt.DeliveryId,
		WaybillID:  opt.WaybillId,
	}

	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}

	return canceler.Cancel(token)
}
