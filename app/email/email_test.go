package email

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/gomail.v2"
	"html/template"
	"testing"
)


func Test_EmailSend(t *testing.T) {

	temlate, err := template.ParseFiles("email_template_table.html")
	if err != nil {
		panic(err)
	}
	var body bytes.Buffer
	m := gomail.NewMessage()
	m.SetHeader("From", "290621352@qq.com")
	m.SetHeader("To", "zhouchangqaz@gmail.com")
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "New Order Notify!")

	temlate.Execute(&body, struct {
		Title string
	}{
		Title: "小周",
	})
	m.SetBody("text/html", body.String())

	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.qq.com", 587, "290621352@qq.com", "sfvtydrkpjscbidb")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}

func TestEmail_Send(t *testing.T) {
	orderId:="12312312312312"
	spew.Dump([]byte(orderId))
	spew.Dump(string([]byte(orderId)))
}