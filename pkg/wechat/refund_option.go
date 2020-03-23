package wechat

import (
	"github.com/iGoogle-ink/gopay"
	wxpay "github.com/iGoogle-ink/gopay/wechat"
	err2 "go-shop-v2/pkg/err"
)

type RefundOption struct {
	appId          string 		// 小程序ID
	mchId          string 		// 商户号
	signType       string 		// 签名类型 MD5,HMAC-SHA256
	nonceStr       string
	OutTradeNo     string		// 商户系统内部订单号
	OutRefundNo    string		// 商户系统内部退款单号
	TotalFee       uint64		// 订单总金额，单位为分
	RefundFee      uint64		// 退款金额，单位分
	RefundDesc     string		// 退款原因
	notifyUrl      string		// 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数
	certFilePath   string
	keyFilePath    string
	pkcs12FilePath string
}

func (refund *RefundOption) SignTypeMD5() *RefundOption {
	refund.signType = wxpay.SignType_MD5
	return refund
}

func (refund *RefundOption) init() {
	if refund.appId == "" {
		refund.appId = Pay.config.AppId
	}
	if refund.mchId == "" {
		refund.mchId = Pay.config.MchId
	}
	if refund.signType == "" {
		refund.SignTypeMD5()
	}
	if refund.notifyUrl == "" {
		refund.notifyUrl = Pay.config.RefundNotifyUrl
	}
	if refund.nonceStr == "" {
		refund.nonceStr = gopay.GetRandomString(32)
	}
	if refund.certFilePath == "" {
		refund.certFilePath = Pay.config.CertFilePath
	}
	if refund.keyFilePath == "" {
		refund.keyFilePath = Pay.config.KeyFilePath
	}
	if refund.pkcs12FilePath == "" {
		refund.pkcs12FilePath = Pay.config.Pkcs12FilePath
	}
}

func (refund *RefundOption) validate() error {
	if refund.appId == "" {
		return err2.Err422.F("发起退款失败，缺少appId")
	}
	if refund.mchId == "" {
		return err2.Err422.F("发起退款失败，缺少mchId")
	}
	if refund.OutTradeNo == "" {
		return err2.Err422.F("发起退款失败，缺少订单号")
	}
	if refund.OutRefundNo == "" {
		return err2.Err422.F("发起退款失败，缺少退款单号")
	}
	if refund.TotalFee == 0 {
		return err2.Err422.F("发起退款失败，订单金额异常")
	}
	if refund.RefundFee == 0 {
		return err2.Err422.F("发起退款失败，退款金额异常")
	}
	if refund.notifyUrl == "" {
		return err2.Err422.F("发起退款失败，缺少回调url")
	}
	if refund.certFilePath == "" {
		return err2.Err422.F("发起退款失败，缺少cert文件地址")
	}
	if refund.keyFilePath == "" {
		return err2.Err422.F("发起退款失败，缺少key文件地址")
	}
	if refund.pkcs12FilePath == "" {
		return err2.Err422.F("发起退款失败，缺少pkcs12文件地址")
	}
	return nil
}

func (refund *RefundOption) toMap() (gopay.BodyMap, error) {
	refund.init()
	if err := refund.validate(); err != nil {
		return nil, err
	}

	bm := make(gopay.BodyMap)
	bm.Set("appid", refund.appId)
	bm.Set("mch_id", refund.mchId)
	bm.Set("nonce_str", refund.nonceStr)
	bm.Set("sign_type", refund.signType)
	bm.Set("out_trade_no", refund.OutTradeNo)
	bm.Set("out_refund_no", refund.OutRefundNo)
	bm.Set("total_fee", refund.TotalFee)
	bm.Set("refund_fee", refund.RefundFee)
	bm.Set("refund_desc", refund.RefundDesc)
	bm.Set("notify_url", refund.notifyUrl)
	// sign
	sign := wxpay.GetParamSign(refund.appId, refund.mchId, Pay.ApiKey, bm)
	// sign, _ := wechat.GetSanBoxParamSign("wxdaa2ab9ef87b5497", mchId, apiKey, body)
	bm.Set("sign", sign)

	return bm, nil
}
