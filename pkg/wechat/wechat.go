package wechat

import (
	"github.com/medivhzhan/weapp/v2"
	"sync"
)

type sdk struct {
	config Config
}

var SDK *sdk
var once sync.Once

func NewSDK(config Config) *sdk {
	once.Do(func() {
		SDK = &sdk{config: config}
	})
	return SDK
}

func (this *sdk) Login(code string) (*weapp.LoginResponse, error) {
	res, err := weapp.Login(this.config.AppId, this.config.AppSecret, code)
	if err != nil {
		// 处理一般错误信息
		return nil, err
	}

	if err := res.GetResponseError(); err != nil {
		// 处理微信返回错误信息
		return nil, err
	}

	return res, nil
	//fmt.Printf("返回结果: %#v", res)
}

func (this *sdk) DecryptUserInfo(sessionKey, rawData, encryptedData, signature, iv string) (userInfo *weapp.UserInfo, err error) {
	return weapp.DecryptUserInfo(sessionKey, rawData, encryptedData, signature, iv)
}
