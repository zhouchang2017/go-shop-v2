package email

import (
	"bytes"
	"go-shop-v2/app/models"
	"html/template"
	"os"
	"path"
)

type orderRefundClosedNotify struct {
	order    *models.Order
	refundId string
	subject  string
	to       string
}

func OrderRefundClosedNotify(order *models.Order, refundId string, to string) *orderRefundClosedNotify {
	return &orderRefundClosedNotify{order: order, refundId: refundId, to: to, subject: "买家关闭退款通知"}
}

func (o orderRefundClosedNotify) To() string {
	return o.to
}

func (o orderRefundClosedNotify) Subject() string {
	return o.subject
}

func (o orderRefundClosedNotify) initData() *refund {
	ref, _ := o.order.FindRefund(o.refundId)
	for _, item := range ref.Items {
		findItem := o.order.FindItem(item.ItemId)
		if findItem != nil {
			item.Item = findItem.Item
		}
	}
	return &refund{Title: "买家关闭退款通知", Refund: ref, OrderId: o.order.GetID()}
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
