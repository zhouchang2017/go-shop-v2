package services

import (
	"context"
	"errors"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/wechat"
	"go-shop-v2/app/models"
	"go-shop-v2/app/repositories"
	"go-shop-v2/config"
	"strconv"
	"time"
)

type PaymentService struct {
	paymentRep *repositories.PaymentRep
	orderResp *repositories.OrderRep
}

func NewPaymentService(paymentRep *repositories.PaymentRep, orderResp *repositories.OrderRep) *PaymentService {
	return &PaymentService{paymentRep: paymentRep, orderResp: orderResp}
}

type PaymentOption struct {
	OrderId string `json:"order_id"`
	OrderNo string `json:"order_no"`
}

func (paymentOpt *PaymentOption) IsValid() error {
	if paymentOpt.OrderId == "" {
		return errors.New("OrderId is empty")
	}
	if paymentOpt.OrderNo == "" {
		return errors.New("OrderNo is empty")
	}
	return nil
}

type WechatMiniPayConfig struct {
	TimeStamp string `json:"timeStamp"`
	NonceStr string `json:"nonceStr"`
	Package string `json:"package"`
	SignType string `json:"signType"`
	PaySign string `json:"paySign"`
}

// 下单
func (srv *PaymentService) Payment(ctx context.Context, userInfo *models.User, opt *PaymentOption) (*WechatMiniPayConfig, error) {
	if err := opt.IsValid(); err != nil {
		return nil, err
	}
	// get order information
	orderRes := <-srv.orderResp.FindById(ctx, opt.OrderId)
	if orderRes.Error != nil {
		return nil, orderRes.Error
	}
	order := orderRes.Result.(*models.Order)
	if order.OrderNo != opt.OrderNo {
		return nil, errors.New("invalid OrderNo with different")
	}
	if order.User.Id != userInfo.ID.String() {
		return nil, errors.New("invalid permission caused of different user")
	}
	if order.Status != models.OrderStatusPrePay {
		return nil, errors.New("invalid order status not prepay")
	}
	// start to generate prepay sign with miniprogram of wechat
	wcPayConfig := config.Config.WechatPayConfig()
	wxClient := wechat.NewClient(wcPayConfig.AppId, wcPayConfig.MchId, wcPayConfig.AppKey, false)
	// unify order
	bodyMap := srv.generateWcBodyMap(order.ActualAmount, order.User.WechatMiniId, order.OrderNo, wechat.SignType_MD5, wcPayConfig.NotifyUrl)
	wxRsp, wcErr := wxClient.UnifiedOrder(bodyMap)
	if wcErr != nil {
		return nil, wcErr
	}
	// prepare params
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	prepayId := "prepay_id" + wxRsp.PrepayId
	paySign := wechat.GetMiniPaySign(wcPayConfig.AppId, wxRsp.NonceStr, prepayId, wechat.SignType_MD5, timeStamp, wcPayConfig.AppKey)
	// update payment information in db
	// todo
	// return
	return &WechatMiniPayConfig{
		TimeStamp: timeStamp,
		NonceStr:  wxRsp.NonceStr,
		Package:   prepayId,
		SignType:  wechat.SignType_MD5,
		PaySign:   paySign,
	}, nil
}

func (srv *PaymentService) generateWcBodyMap(amount uint64, openId, orderNo, signType, notifyUrl string) gopay.BodyMap {
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", gopay.GetRandomString(32))
	bm.Set("body", "小程序测试支付")			//todo
	bm.Set("out_trade_no", orderNo)
	bm.Set("total_fee", amount)
	bm.Set("spbill_create_ip", "127.0.0.1")		//todo
	bm.Set("notify_url", notifyUrl)
	bm.Set("trade_type", wechat.TradeType_Mini)
	bm.Set("device_info", "WEB")
	bm.Set("sign_type", signType)
	bm.Set("openid", openId)

	//// 嵌套json格式数据（例如：H5支付的 scene_info 参数）
	//h5Info := make(map[string]string)
	//h5Info["type"] = "Wap"
	//h5Info["wap_url"] = "http://www.gopay.ink"
	//h5Info["wap_name"] = "H5测试支付"
	//
	//sceneInfo := make(map[string]map[string]string)
	//sceneInfo["h5_info"] = h5Info
	//bm.Set("scene_info", sceneInfo)
	// return
	return bm
}