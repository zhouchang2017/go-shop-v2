package email

import (
	"fmt"
	"go-shop-v2/pkg/utils"
	"gopkg.in/gomail.v2"
	"strconv"
	"sync"
	"time"
)

type Mailer interface {
	Send(to string, subject string, body string) error
}
type Notify interface {
	Subject() string                // 主题
	Body() (body string, err error) // 内容
}

type NotifyContentType interface {
	ContentType() string // 内容类型
}

type Receiver interface {
	GetNickname() string
	GetEmail() string
}

var instance *email
var once sync.Once

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Sender   string `json:"sender"`
}

func New(c Config) *email {
	once.Do(func() {
		instance = &email{config: c}
		instance.init()
	})
	return instance
}

type email struct {
	config Config
	*gomail.Dialer
	projectPath string
}

func (e *email) init() {
	e.Dialer = gomail.NewDialer(e.config.Host, e.config.Port, e.config.Username, e.config.Password)
}

func (e *email) Send(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", instance.config.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return e.DialAndSend(m)
}

func Send(notify Notify, to string, cc ...Receiver) error {
	m := gomail.NewMessage()
	m.SetHeader("From", instance.config.Username)
	m.SetHeader("To", to)
	for _, c := range cc {
		m.SetAddressHeader("Cc", c.GetEmail(), c.GetNickname())
	}

	m.SetHeader("Subject", notify.Subject())

	body, err := notify.Body()
	if err != nil {
		return err
	}
	if customType, ok := notify.(NotifyContentType); ok {
		m.SetBody(customType.ContentType(), body)
	} else {
		m.SetBody("text/html", body)
	}

	return instance.DialAndSend(m)
}

func Sends(notify Notify, receivers ...Receiver) error {
	to := receivers[0]

	if len(receivers) > 1 {
		cc := receivers[1:]
		return Send(notify, to.GetEmail(), cc...)
	}
	return Send(notify, to.GetEmail())
}

// 过滤器
// 时间转换
func timeStr(t time.Time) string {
	return utils.TimeJsonOut(t)
}

// 金额转换
func money(amount interface{}) string {
	var price float64
	switch amount.(type) {
	case int64:
		price = float64(amount.(int64)) / 100
	case uint64:
		price = float64(amount.(uint64)) / 100
	case int:
		price = float64(amount.(int)) / 100
	default:
		price = 0
	}
	float := strconv.FormatFloat(price, 'f', 2, 64)
	return fmt.Sprintf("￥%s", float)
}
