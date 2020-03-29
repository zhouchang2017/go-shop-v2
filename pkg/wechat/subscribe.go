package wechat

import (
	"encoding/json"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/medivhzhan/weapp/v2"
	"log"
	"net/http"
	"net/url"
	"reflect"
)

type Sender interface {
	To() string
	TemplateID() string
	Page() string
	Data() weapp.SubscribeMessageData
}

// 订阅消息
func (this *sdk) SendSubscribeMessage(sender Sender) error {

	token, err := this.getAccessToken()
	if err != nil {
		return err
	}
	message := weapp.SubscribeMessage{
		ToUser:     sender.To(),
		TemplateID: sender.TemplateID(),
		Data:       sender.Data(),
	}
	if sender.Page() != "" {
		message.Page = sender.Page()
	}

	spew.Dump("发送订阅消息数据", message)

	send, err := message.Send(token)
	if err != nil {
		return err
	}
	if err := send.GetResponseError(); err != nil {
		log.Printf("SendSubscribeMessage Error,errCode:%d,errMsg:%s", send.ErrCode, send.ErrMSG)
		if send.ErrCode == 43101 {
			// 用户未开启
			return nil
		}
		return err
	}
	return nil
}

type TemplateInfoResponse struct {
	weapp.CommonError
	Count int                 `json:"count"` // 模版标题列表总数
	Data  []*TemplateInfoData `json:"data"`  // 关键词列表
}

type TemplateInfoData struct {
	Kid     int    `json:"kid"`     // 关键词 id，选用模板时需要
	Name    string `json:"name"`    // 关键词内容
	Example string `json:"example"` // 关键词内容对应的示例
	Rule    string `json:"rule"`    //参数类型
}

// 获取模板标题下的关键词列表
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getPubTemplateKeyWordsById.html
func (this *sdk) GetPubTemplateKeyWordsById(tid string) (res *TemplateInfoResponse, err error) {
	const uri = "https://api.weixin.qq.com/wxaapi/newtmpl/getpubtemplatekeywords"
	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}
	build, err := this.build(uri, map[string]interface{}{
		"access_token": token,
		"tid":          tid,
	})
	if err != nil {
		// uri构建错误
		return nil, err
	}
	response, err := http.Get(build)
	if err != nil {
		// http 请求错误
		return nil, err
	}
	res = &TemplateInfoResponse{}
	if err := json.NewDecoder(response.Body).Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}

type TemplateListResponse struct {
	weapp.CommonError
	Data []*TemplateListData `json:"data"`
}

type TemplateListData struct {
	PriTmplId string `json:"priTmplId"` // 添加至帐号下的模板 id，发送小程序订阅消息时所需
	Title     string `json:"title"`     // 模版标题
	Content   string `json:"content"`   // 模版内容
	Example   string `json:"example"`   // 模板内容示例
	Type      int    `json:"type"`      // 模版类型，2 为一次性订阅，3 为长期订阅
}

// 获取当前帐号下的个人模板列表
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getTemplateList.html
func (this *sdk) GetTemplateList() (res *TemplateListResponse, err error) {
	const uri = "https://api.weixin.qq.com/wxaapi/newtmpl/gettemplate"
	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}
	build, err := this.build(uri, map[string]interface{}{
		"access_token": token,
	})
	if err != nil {
		// uri构建错误
		return nil, err
	}
	response, err := http.Get(build)
	if err != nil {
		// http 请求错误
		return nil, err
	}
	res = &TemplateListResponse{}
	if err := json.NewDecoder(response.Body).Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}
func (this *sdk) build(uri string, query map[string]interface{}) (string, error) {
	api, _ := url.Parse(uri)
	q := url.Values{}
	for key, value := range query {
		typeOf := reflect.TypeOf(value)
		var v string
		switch typeOf.Kind() {
		case reflect.String:
			v = value.(string)
		case reflect.Map:
			bytes, err := json.Marshal(value)
			if err != nil {
				return "", err
			}
			v = string(bytes)
		default:
			return "", errors.New("value 必须为 string 或者 map")
		}
		q.Add(key, v)
	}
	api.RawQuery = q.Encode()

	return api.String(), nil
}
