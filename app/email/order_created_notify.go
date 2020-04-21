package email

import (
	"bytes"
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

func (this *orderCreatedNotify) Body() (string, error) {
	getwd, _ := os.Getwd()
	fileName := "order_created_notify.html"
	filePath := path.Join(getwd, "app", "email", "template", fileName)
	content, err := template.New(fileName).Funcs(this.funcMap()).ParseFiles(filePath)
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

func OrderCreatedNotify(order *models.Order) *orderCreatedNotify {
	return &orderCreatedNotify{
		order:   order,
		subject: "新订单通知",
	}
}

func (this *orderCreatedNotify) Subject() string {
	return this.subject
}

func (this *orderCreatedNotify) funcMap() template.FuncMap {
	return template.FuncMap{
		"timeStr": timeStr,
		"money":   money,
	}
}
