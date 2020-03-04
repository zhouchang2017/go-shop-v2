package lbs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sync"
)

const URI = "https://apis.map.qq.com/ws/geocoder/v1/"

var SDK *sdk

type sdk struct {
	key string
}

var once sync.Once

func NewSDK(key string) *sdk {
	once.Do(func() {
		SDK = &sdk{key: key}
	})
	return SDK
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

type parseJsonStruct struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Location Location `json:"location"`
	} `json:"result"`
}

type Location struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

func (this sdk) ParseAddress(addr string, options ...map[string]string) (res *Location, err error) {
	res = &Location{}
	option := map[string]interface{}{}
	for _, opt := range options {
		for key, value := range opt {
			option[key] = value
		}
	}

	option["address"] = addr
	option["key"] = this.key

	uri, err := this.build(URI, option)

	resp, err := http.Get(uri)
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var parse parseJsonStruct
	if err = json.Unmarshal(bytes, &parse); err != nil {
		return
	}

	if parse.Status == 0 {
		return &parse.Result.Location, nil
	}
	return res, errors.New("lbs skd 异常")
}
