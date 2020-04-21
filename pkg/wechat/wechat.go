package wechat

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/medivhzhan/weapp/v2"
	"github.com/nfnt/resize"
	"go-shop-v2/pkg/cache/redis"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
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

// 是否为生产环境
func isProd() bool {
	return os.Getenv("GIN_MODE") == "release"
}

// 用户登录
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

func (this *sdk) IsProd() bool {
	return this.config.IsProd
}

// 解密用户信息
func (this *sdk) DecryptUserInfo(sessionKey, rawData, encryptedData, signature, iv string) (userInfo *weapp.UserInfo, err error) {
	return weapp.DecryptUserInfo(sessionKey, rawData, encryptedData, signature, iv)
}

const accessTokenCacheKey = "wechat_access_token"

func ClearCache() {
	if redis.GetConFn() != nil {
		redis.GetConFn().Del(accessTokenCacheKey)
	}
}

var testAccessToken string

// get access token
func (this *sdk) getAccessToken() (accessToken string, err error) {
	if testAccessToken != "" {
		return testAccessToken, nil
	}
	if redis.GetConFn() != nil {
		result, err := redis.GetConFn().Get(accessTokenCacheKey).Result()

		if err == nil {
			log.Printf("get access token from cache,token = %s", result)
			return result, nil
		}

	}
	tryNum := 3
GETTOKEN:
	token, err := weapp.GetAccessToken(this.config.AppId, this.config.AppSecret)
	if err != nil {
		return "", err
	}

	if token.ErrCode != 0 {
		// 调用失败
		if token.ErrCode == -1 && tryNum > 0 {
			// 系统繁忙，此时请开发者稍候再试
			after := time.After(time.Millisecond * 500)
			<-after
			tryNum--
			if tryNum == 0 {
				err = fmt.Errorf("Wechat GetAccessToken Error,try 3,err:%s", token.ErrMSG)
				return
			}
			goto GETTOKEN
		}
	}
	if redis.GetConFn() != nil {
		expires := time.Duration(token.ExpiresIn) - 200
		redis.GetConFn().Set(accessTokenCacheKey, token.AccessToken, time.Second*expires)
	}
	return token.AccessToken, nil
}

func timeFormat(t time.Time) string {
	return t.Format("20060102")
}

// 访问留存
// 获取用户访问小程序日留存
func (this *sdk) GetDailyRetain(start, end time.Time) (*weapp.RetainResponse, error) {
	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}
	retain, err := weapp.GetDailyRetain(token, timeFormat(start), timeFormat(end))
	if err != nil {
		return nil, err
	}
	if err := retain.GetResponseError(); err != nil {
		return nil, err
	}
	return retain, nil
}

// 获取用户访问小程序数据概况
// getDailySummary
func (this *sdk) GetDailySummary(day time.Time) (*weapp.DailySummary, error) {
	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}
	end := day.AddDate(0, 0, -1)
	summary, err := weapp.GetDailySummary(token, timeFormat(end), timeFormat(end))
	if err != nil {
		return nil, err
	}

	if err := summary.GetResponseError(); err != nil {
		return nil, err
	}

	return summary, nil
}

// 获取用户访问小程序数据日趋势
// getDailyVisitTrend
func (this *sdk) GetDailyVisitTrend(day time.Time) (*weapp.VisitTrend, error) {
	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}
	endTime := time.Now().AddDate(0, 0, -1)

	if day.After(endTime) {
		day = endTime
	}

	trend, err := weapp.GetDailyVisitTrend(token, timeFormat(day), timeFormat(day))
	if err != nil {
		return nil, err
	}
	if err := trend.GetResponseError(); err != nil {
		return nil, err
	}

	return trend, nil
}

func (this *sdk) NewServer() (*weapp.Server, error) {
	return weapp.NewServer(this.config.AppId, this.config.Token, this.config.EncodingAESKey, Pay.MchId, Pay.ApiKey, false)
}

// 生成小程序码
func (this *sdk) UnlimitedQRCode(opt weapp.UnlimitedQRCode) ([]byte, error) {
	if !this.IsProd() {
		projectName := "go-shop-v2"
		configFileName := "static/img/gh_c764313ca22e_344.png"
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		split := strings.Split(dir, projectName)

		filePath := path.Join(split[0], projectName, configFileName)
		open, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer open.Close()
		img, _, err := image.Decode(open)
		if err != nil {
			return nil, err
		}
		var width uint
		if opt.Width == 0 {
			width = 280
		} else {
			width = uint(opt.Width)
		}
		resizeImg := resize.Resize(width, width, img, resize.Lanczos3)
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, resizeImg); err != nil {
			return nil, err
		}
		base64.StdEncoding.EncodeToString(buf.Bytes())
		return buf.Bytes(), nil
	}
	token, err := this.getAccessToken()
	if err != nil {
		return nil, err
	}
	resp, res, err := opt.Get(token)
	if err != nil {
		// 处理一般错误信息
		return nil, err
	}

	if err := res.GetResponseError(); err != nil {
		// 处理微信返回错误信息
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
