package email

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/models"
	"html/template"
	"os"
	"path"
)

type orderCreatedNotify struct {
	order   *models.Order
	subject string
	to      string
}

func (this *orderCreatedNotify) To() string {
	return this.to
}

func (this *orderCreatedNotify) Body() (string, error) {
	getwd, _ := os.Getwd()
	filePath := path.Join(getwd, "app", "email", "template", "order_created_notify.html")
	spew.Dump(filePath)
	content, err := template.New("order_created_notify.html").Funcs(this.funcMap()).ParseFiles(filePath)
	if err != nil {
		return "", err
	}
	//content.Funcs(funcMap)

	var body bytes.Buffer

	if err := content.Execute(&body, this.order); err != nil {
		return "", err
	}

	return body.String(), nil
}

func OrderCreatedNotify(order *models.Order, to string) *orderCreatedNotify {
	return &orderCreatedNotify{
		order:   order,
		subject: "新订单通知",
		to:      to,
	}
}

func (this *orderCreatedNotify) Subject() string {
	return this.subject
}

func (this *orderCreatedNotify) funcMap() template.FuncMap {
	return template.FuncMap{
		// 注册函数title, strings.Title会将单词首字母大写
		"timeStr": timeStr,
		"money":   money,
	}
}
