package email

import (
	"bytes"
	"go-shop-v2/app/models"
	"html/template"
	"os"
	"path"
)

// 买家关闭订单通知
type orderClosedNotify struct {
	order   *models.Order
	subject string
	to      string
}

func OrderClosedNotify(order *models.Order) *orderClosedNotify {
	return &orderClosedNotify{order: order,  subject: "买家取消订单通知",}
}

func (o orderClosedNotify) Subject() string {
	return o.subject
}

func (o orderClosedNotify) Body() (string, error) {
	getwd, _ := os.Getwd()
	fileName := "order_closed_notify.html"
	filePath := path.Join(getwd, "app", "email", "template", fileName)
	content, err := template.New(fileName).Funcs(o.funcMap()).ParseFiles(filePath)
	if err != nil {
		return "", err
	}
	//content.Funcs(funcMap)

	var body bytes.Buffer

	if err := content.Execute(&body, o.order); err != nil {
		return "", err
	}

	return body.String(), nil
}

func (this *orderClosedNotify) funcMap() template.FuncMap {
	return template.FuncMap{
		"timeStr": timeStr,
		"money":   money,
	}
}
