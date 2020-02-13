package qiniu

import (
	"context"
	"github.com/gin-gonic/gin"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/utils"
	"net/http"
)

var QiniuController = qiniuController{}

type qiniuController struct {
}

func (qiniuController) Token(c *gin.Context) {
	token, err := instance.Token(context.Background())
	if err != nil {
		err2.ErrorEncoder(nil, err, c.Writer)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// 第三方链接转存七牛云
func (qiniuController) Fetch(c *gin.Context) {
	uri := c.PostForm("url")
	if uri == "" {
		err2.ErrorEncoder(nil, err2.NewFromCode(422).F("url参数不能为空"), c.Writer)
		return
	}
	res, err := instance.PutByUrl(c, uri, utils.RandomString(32))
	if err != nil {
		err2.ErrorEncoder(nil, err, c.Writer)
		return
	}
	c.JSON(http.StatusOK, res)
}

// 富文本中第三方链接转存七牛云
func (qiniuController) RichTextFetch(c *gin.Context) {

}
