package email

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/app/models"
	"html/template"
	"os"
	"path"
)

type orderPaidNotify struct {
	order   *models.Order
	subject string
	to      string
}

func (this *orderPaidNotify) To() string {
	return this.to
}

func (this *orderPaidNotify) Body() (string, error) {
	getwd, _ := os.Getwd()
	filePath := path.Join(getwd, "app", "email", "template", "order_paid_notify.html")
	spew.Dump(filePath)
	content, err := template.New("order_paid_notify.html").Funcs(this.funcMap()).ParseFiles(filePath)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	if err := content.Execute(&body, this.order); err != nil {
		return "", err
	}

	return body.String(), nil
}

func OrderPaidNotify(order *models.Order, to string) *orderPaidNotify {
	return &orderPaidNotify{
		order:   order,
		subject: "订单付款通知",
		to:      to,
	}
}

func (this *orderPaidNotify) Subject() string {
	return this.subject
}

func (this *orderPaidNotify) funcMap() template.FuncMap {
	return template.FuncMap{
		// 注册函数title, strings.Title会将单词首字母大写
		"timeStr": timeStr,
		"money":   money,
	}
}
