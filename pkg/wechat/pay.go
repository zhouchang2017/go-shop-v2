package wechat

import (
	"fmt"
	wxpay "github.com/iGoogle-ink/gopay/wechat"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type pay struct {
	config PayConfig
	*wxpay.Client
}

var Pay *pay
var oncePay sync.Once

func NewPay(config PayConfig) *pay {
	oncePay.Do(func() {
		Pay = &pay{config: config}
		client := wxpay.NewClient(config.AppId, config.MchId, config.AppKey, config.IsProd)
		client.SetCountry(wxpay.China)
		client.AddCertFilePath(config.CertFilePath, config.KeyFilePath, config.Pkcs12FilePath)
		Pay.Client = client
	})
	return Pay
}

type unifiedOrderResponse struct {
	*wxpay.UnifiedOrderResponse
	timestamp string
	signType  string
}

type WechatMiniPayConfig struct {
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

func (u *unifiedOrderResponse) GetWechatMiniPayConfig() *WechatMiniPayConfig {
	paySign := u.GetMiniPaySign()
	return &WechatMiniPayConfig{
		TimeStamp: u.timestamp,
		NonceStr:  u.NonceStr,
		Package:   fmt.Sprintf("prepay_id=%s", u.PrepayId),
		SignType:  u.signType,
		PaySign:   paySign,
	}
}

func (u *unifiedOrderResponse) setTimeStamp() {
	u.timestamp = strconv.FormatInt(time.Now().Unix(), 10)
}

// ====APP支付 paySign====
func (u *unifiedOrderResponse) GetAppPaySign(partnerid string) string {
	u.setTimeStamp()
	// 获取APP支付的 paySign
	// 注意：package 参数因为是固定值，无需开发者再传入
	//    appId：AppID
	//    partnerid：partnerid
	//    nonceStr：随机字符串
	//    prepayId：统一下单成功后得到的值
	//    signType：签名方式，务必与统一下单时用的签名方式一致
	//    timeStamp：时间
	//    apiKey：API秘钥值
	return wxpay.GetAppPaySign(Pay.AppId, partnerid, u.NonceStr, u.PrepayId, u.signType, u.timestamp, Pay.ApiKey)
}

// ====微信小程序 paySign====
func (u *unifiedOrderResponse) GetMiniPaySign() string {
	u.setTimeStamp()
	prepayId := fmt.Sprintf("prepay_id=%s", u.PrepayId) // 此处的 wxRsp.PrepayId ,统一下单成功后得到
	// 获取微信小程序支付的 paySign
	//    appId：AppID
	//    nonceStr：随机字符串
	//    prepayId：统一下单成功后得到的值
	//    signType：签名方式，务必与统一下单时用的签名方式一致
	//    timeStamp：时间
	//    apiKey：API秘钥值

	return wxpay.GetMiniPaySign(Pay.AppId, u.NonceStr, prepayId, u.signType, u.timestamp, Pay.ApiKey)
}

// ====微信内H5支付 paySign====
func (u *unifiedOrderResponse) GetH5PaySign() string {
	u.setTimeStamp()
	packages := fmt.Sprintf("prepay_id=%s", u.PrepayId) // 此处的 wxRsp.PrepayId ,统一下单成功后得到
	// 获取微信内H5支付 paySign
	//    appId：AppID
	//    nonceStr：随机字符串
	//    packages：统一下单成功后拼接得到的值
	//    signType：签名方式，务必与统一下单时用的签名方式一致
	//    timeStamp：时间
	//    apiKey：API秘钥值
	return wxpay.GetH5PaySign(Pay.AppId, u.NonceStr, packages, u.signType, u.timestamp, Pay.ApiKey)
}

// 统一下单
func (this pay) UnifiedOrder(opt *PayUnifiedOrderOption) (res *unifiedOrderResponse, err error) {
	bm, err := opt.toMap()
	if err != nil {
		return nil, err
	}
	wxRsp, err := this.Client.UnifiedOrder(bm)
	if err != nil {
		return nil, err
	}
	_, err = wxpay.VerifySign(this.ApiKey, opt.signType, wxRsp)
	if err != nil {
		return nil, err
	}
	return &unifiedOrderResponse{UnifiedOrderResponse: wxRsp, signType: opt.signType}, nil
}

// 支付异步通知参数解析和验签Sign
func (this pay) ParseNotifyResult(req *http.Request) (notifyReq *wxpay.NotifyRequest, err error) {
	notifyReq, err = wxpay.ParseNotifyResult(req)
	if err != nil {
		return
	}
	_, err = wxpay.VerifySign(this.ApiKey, wxpay.SignType_MD5, notifyReq)
	return
}
