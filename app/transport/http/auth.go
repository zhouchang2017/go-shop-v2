package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/app/services"
	"go-shop-v2/pkg/auth"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/wechat"
	"net/http"
)

type AuthController struct {
	userSrv *services.UserService
}

// 用户登陆
func (this *AuthController) Login(ctx *gin.Context) {
	var form loginForm
	if err := ctx.ShouldBind(&form); err != nil {
		ResponseError(ctx, err)
	}
	if form.Code == "" {
		ResponseError(ctx, err2.Error(401, "code无效"))
		return
	}
	res, err := wechat.SDK.Login(form.Code)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	jwtGuard, err := auth.Auth.Guard(guard)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	credentials := map[string]string{
		"open_id": res.OpenID,
	}

	if res, ok := jwtGuard.Attempt(credentials, true); ok {
		// todo 更新用户数据
		token := fmt.Sprintf("%s", res)
		ctx.Header("token", token)
		Response(ctx, gin.H{
			"code":  200,
			"token": token,
		}, http.StatusOK)
		return
	}
	// 注册
	info, err := wechat.SDK.DecryptUserInfo(res.SessionKey, form.RawData, form.EncryptedData, form.Signature, form.Iv)
	if err !=nil {
		ResponseError(ctx, err)
		return
	}

	// 注册新用户
	user, err := this.userSrv.RegisterByWechat(ctx, info)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	data, err := jwtGuard.Login(user)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	token := fmt.Sprintf("%s", data)
	ctx.Header("token", token)
	Response(ctx, gin.H{
		"code":  200,
		"token": token,
	}, http.StatusOK)

}

type loginForm struct {
	Code          string `json:"code"`
	RawData       string `json:"rawData"`
	Signature     string `json:"signature"`
	Iv            string `json:"iv"`
	EncryptedData string `json:"encryptedData"`
}