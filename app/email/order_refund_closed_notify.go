package email

import (
	"bytes"
	"go-shop-v2/app/models"
	"html/template"
	"os"
	"path"
)

type orderRefundClosedNotify struct {
	refund  *models.Refund
	subject string
}

func OrderRefundClosedNotify(refund *models.Refund) *orderRefundClosedNotify {
	return &orderRefundClosedNotify{refund: refund, subject: "买家关闭退款通知"}
}

func (o orderRefundClosedNotify) Subject() string {
	return o.subject
}

func (o orderRefundClosedNotify) initData() *refund {
	return &refund{Title: "买家关闭退款通知", Refund: o.refund, OrderId: o.refund.OrderId}
}

func (o orderRefundClosedNotify) Body() (string, error) {
	getwd, _ := os.Getwd()
	fileName := "order_apply_refund_notify.html"
	filePath := path.Join(getwd, "app", "email", "template", fileName)
	content, err := template.New(fileName).Funcs(o.funcMap()).ParseFiles(filePath)
	if err != nil {
		return "", err
	}
	//content.Funcs(funcMap)

	var body bytes.Buffer

	if err := content.Execute(&body, o.initData()); err != nil {
		return "", err
	}

	return body.String(), nil
}

func (this *orderRefundClosedNotify) funcMap() template.FuncMap {
	return template.FuncMap{
		"timeStr": timeStr,
		"money":   money,
	}
}
