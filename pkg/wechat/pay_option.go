package wechat

import (
	"encoding/json"
	"github.com/iGoogle-ink/gopay"
	wxpay "github.com/iGoogle-ink/gopay/wechat"
	err2 "go-shop-v2/pkg/err"
	"log"
	"time"
)

// 单品优惠活动detail字段列表
type PayOrderDetail struct {
	CostPrice   uint64                 `json:"cost_price"` // 订单原价
	GoodsDetail []*PayOrderGoodsDetail `json:"goods_detail"`
}

func (p PayOrderDetail) toJsonString() (string, error) {
	marshal, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// 单品优惠活动goods_detail字段
type PayOrderGoodsDetail struct {
	GoodsId   string `json:"goods_id"`   // 商品编码
	GoodsName string `json:"goods_name"` // 商品名称
	Quantity  uint64  `json:"quantity"`   // 数量
	Price     uint64  `json:"price"`      // 如果商户有优惠，需传输商户优惠后的单价
}

type PayUnifiedOrderOption struct {
	appId          string          // 小程序ID
	mchId          string          // 商户号
	DeviceInfo     string          // 自定义参数，可以为终端设备号(门店号或收银设备ID)，PC网页或公众号内支付可以传"WEB"
	signType       string          // 签名类型 MD5,HMAC-SHA256
	Body           string          // 商品描述
	Detail         *PayOrderDetail // 商品详细描述，对于使用单品优惠的商户，该字段必须按照规范上传，详见“单品优惠参数说明”
	attach         string          // 附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。
	OutTradeNo     string          // 商户系统内部订单号
	TotalFee       uint64          // 订单总金额，单位为分
	SpbillCreateIp string          // 支持IPV4和IPV6两种格式的IP地址。调用微信支付API的机器IP
	timeStart      string          // 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。
	timeExpire     string          // 订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。订单失效时间是针对订单号而言的，由于在请求支付的时候有一个必传参数prepay_id只有两小时的有效期，所以在重入时间超过2小时的时候需要重新请求下单接口获取新的prepay_id。
	goodsTag       string          // 订单优惠标记，使用代金券或立减优惠功能时需要的参数
	notifyUrl      string          // 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数
	tradeType      string          // 小程序取值如下：JSAPI
	OpenId         string          // 用户在商户appid下的唯一标识。
	nonceStr       string
}

func (p *PayUnifiedOrderOption) AppId(id string) *PayUnifiedOrderOption {
	p.appId = id
	return p
}

func (p *PayUnifiedOrderOption) MchId(id string) *PayUnifiedOrderOption {
	p.mchId = id
	return p
}

func (p *PayUnifiedOrderOption) SignTypeMD5() *PayUnifiedOrderOption {
	p.signType = wxpay.SignType_MD5
	return p
}

func (p *PayUnifiedOrderOption) SignTypeHMAC_SHA256() *PayUnifiedOrderOption {
	p.signType = wxpay.SignType_HMAC_SHA256
	return p
}

func (p *PayUnifiedOrderOption) Attach(attach string) *PayUnifiedOrderOption {
	p.attach = attach
	return p
}

func (p *PayUnifiedOrderOption) TimeStart(t time.Time) *PayUnifiedOrderOption {
	p.timeStart = t.Format("20060102150405")
	return p
}

func (p *PayUnifiedOrderOption) TimeExpire(t time.Time) *PayUnifiedOrderOption {
	p.timeExpire = t.Format("20060102150405")
	return p
}

func (p *PayUnifiedOrderOption) GoodsTag(tag string) *PayUnifiedOrderOption {
	p.goodsTag = tag
	return p
}

func (p *PayUnifiedOrderOption) NotifyUrl(url string) *PayUnifiedOrderOption {
	p.notifyUrl = url
	return p
}

func (p *PayUnifiedOrderOption) init() {
	if p.appId == "" {
		p.appId = Pay.config.AppId
	}
	if p.mchId == "" {
		p.mchId = Pay.config.MchId
	}
	if p.DeviceInfo == "" {
		p.DeviceInfo = "MP"
	}
	if p.signType == "" {
		p.SignTypeMD5()
	}
	if p.notifyUrl == "" {
		p.notifyUrl = Pay.config.NotifyUrl
	}
	if p.tradeType == "" {
		p.tradeType = wxpay.TradeType_Mini
	}
	if p.nonceStr == "" {
		p.nonceStr = gopay.GetRandomString(32)
	}
}

func (p *PayUnifiedOrderOption) validate() error {
	if p.appId == "" {
		return err2.Err422.F("发起支付失败，缺少appId")
	}
	if p.mchId == "" {
		return err2.Err422.F("发起支付失败，缺少mchId")
	}
	if p.OutTradeNo == "" {
		return err2.Err422.F("发起支付失败，缺少订单号")
	}
	if p.TotalFee == 0 {
		return err2.Err422.F("发起支付失败，付款金额异常")
	}
	if p.Body == "" {
		return err2.Err422.F("发起支付失败，缺少支付商品描述")
	}
	if p.notifyUrl == "" {
		return err2.Err422.F("发起支付失败，缺少回调url")
	}
	if p.OpenId == "" {
		return err2.Err422.F("发起支付失败，缺少openId")
	}
	if p.notifyUrl == "" {
		err2.Err422.F("发起支付失败，缺少notifyUrl")
	}
	return nil
}

func (p *PayUnifiedOrderOption) toMap() (gopay.BodyMap, error) {
	p.init()
	if err := p.validate(); err != nil {
		return nil, err
	}

	bm := make(gopay.BodyMap)
	bm.Set("appid", p.appId)
	bm.Set("mch_id", p.mchId)
	bm.Set("device_info", p.DeviceInfo)
	bm.Set("nonce_str", p.nonceStr)
	bm.Set("sign_type", p.signType)
	bm.Set("body", p.Body)
	if p.Detail != nil {
		jsonString, err := p.Detail.toJsonString()
		if err == nil {
			bm.Set("detail", jsonString)
		} else {
			log.Printf("p.Detail.toJsonString() error：%s\n", err)
		}
	}
	if p.attach != "" {
		bm.Set("attach", p.attach)
	}
	bm.Set("out_trade_no", p.OutTradeNo)
	bm.Set("total_fee", p.TotalFee)
	bm.Set("spbill_create_ip", p.SpbillCreateIp)
	if p.timeStart != "" && p.timeExpire != "" {
		bm.Set("time_start", p.timeStart)
		bm.Set("time_expire", p.timeExpire)
	}
	if p.goodsTag != "" {
		bm.Set("goods_tag", p.goodsTag)
	}
	bm.Set("notify_url", p.notifyUrl)
	bm.Set("trade_type", p.tradeType)
	bm.Set("openid", p.OpenId)

	return bm, nil
}
