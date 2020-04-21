package email

import (
	"bytes"
	"go-shop-v2/app/models"
	"html/template"
	"os"
	"path"
)

// 买家申请退款通知
type orderApplyRefundNotify struct {
	refund  *models.Refund
	subject string
	to      string
}

func OrderApplyRefundNotify(refund *models.Refund) *orderApplyRefundNotify {
	return &orderApplyRefundNotify{refund: refund, subject: "买家申请退款通知"}
}

func (o orderApplyRefundNotify) Subject() string {
	return o.subject
}

type refund struct {
	Title string
	*models.Refund
	OrderId string
}

func (o orderApplyRefundNotify) initData() *refund {
	return &refund{Title: "买家申请退款通知", Refund: o.refund}
}

func (o orderApplyRefundNotify) Body() (string, error) {
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

func (this *orderApplyRefundNotify) funcMap() template.FuncMap {
	return template.FuncMap{
		"timeStr": timeStr,
		"money":   money,
	}
}
