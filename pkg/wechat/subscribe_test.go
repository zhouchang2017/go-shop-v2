package wechat

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestSdk_GetPubTemplateKeyWordsById(t *testing.T) {

	newSDK := NewSDK(Config{})
	testAccessToken = ""
	res, err := newSDK.GetPubTemplateKeyWordsById("994")
	if err!=nil {
		panic(err)
	}
	spew.Dump(res)

}

func TestSdk_GetTemplateList(t *testing.T) {
	config := Config{
		AppId:     "",
		AppSecret: "",
	}
	newSDK := NewSDK(config)
	list, err := newSDK.GetTemplateList()
	if err!=nil {
		panic(err)
	}
	spew.Dump(list)
}
