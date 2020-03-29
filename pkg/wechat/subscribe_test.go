package wechat

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestSdk_GetPubTemplateKeyWordsById(t *testing.T) {

	newSDK := NewSDK(Config{})
	testAccessToken = "31_jngvPVsUKy6MTqvlV8OZd5x7cH1uf6-Yvjm4A9RXX--fW6IkMIha4y_Nz5mSpXX15aYxoQJdS2nQHTMz4jU-5EWIwWsDpeXNuevOp9phL6thmZT48bxW3urR_JtiaHlDugdH4B8MYdcEmza_GWQcACAQGA"
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
