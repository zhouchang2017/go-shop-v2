package email

import (
	"bytes"
	"go-shop-v2/app/models"
	"html/template"
	"os"
	"path"
)

type orderPaidNotify struct {
	order   *models.Order
	subject string
}

func (this *orderPaidNotify) Body() (string, error) {
	getwd, _ := os.Getwd()
	fileName := "order_paid_notify.html"
	filePath := path.Join(getwd, "app", "email", "template", fileName)
	content, err := template.New(fileName).Funcs(this.funcMap()).ParseFiles(filePath)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	if err := content.Execute(&body, this.order); err != nil {
		return "", err
	}

	return body.String(), nil
}

func OrderPaidNotify(order *models.Order) *orderPaidNotify {
	return &orderPaidNotify{
		order:   order,
		subject: "订单付款通知",
	}
}

func (this *orderPaidNotify) Subject() string {
	return this.subject
}

func (this *orderPaidNotify) funcMap() template.FuncMap {
	return template.FuncMap{
		"timeStr": timeStr,
		"money":   money,
	}
}
