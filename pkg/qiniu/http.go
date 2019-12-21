package qiniu

import (
	"context"
	"github.com/gin-gonic/gin"
	err2 "go-shop-v2/pkg/err"
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
