package email

import (
	"flag"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailParam struct {
	// ServerHost 邮箱服务器地址，如腾讯企业邮箱为smtp.exmail.qq.com
	ServerHost string
	// ServerPort 邮箱服务器端口，如腾讯企业邮箱为465
	ServerPort int
	// FromEmail　发件人邮箱地址
	FromEmail string
	// FromPasswd 发件人邮箱密码（注意，这里是明文形式），TODO：如果设置成密文？
	FromPasswd string
	// Toers 接收者邮件，如有多个，则以英文逗号(“,”)隔开，不能为空
	Toers string
	// CCers 抄送者邮件，如有多个，则以英文逗号(“,”)隔开，可以为空
	CCers string
}

var m *gomail.Message

var serverHost = flag.String("serverHost", "smtp.exmail.qq.com", "email smtp host")
var serverPort = flag.Int("serverPort", 465, "email smtp post")
var fromEmail = flag.String("emailName", "", "email address")
var fromEmailPasswd = flag.String("emailPasswd", "", "email smtp password")

func initEmail(ep *EmailParam) {
	toers := []string{}

	m = gomail.NewMessage()

	if len(ep.Toers) == 0 {
		return
	}

	for _, tmp := range strings.Split(ep.Toers, ",") {
		toers = append(toers, strings.TrimSpace(tmp))
	}

	// 收件人可以有多个，故用此方式
	m.SetHeader("To", toers...)

	//抄送列表
	if len(ep.CCers) != 0 {
		for _, tmp := range strings.Split(ep.CCers, ",") {
			toers = append(toers, strings.TrimSpace(tmp))
		}
		m.SetHeader("Cc", toers...)
	}

	// 发件人
	m.SetAddressHeader("From", ep.FromEmail, "Github Action Hinter")
}

// sendEmail body支持html格式字符串
func sendEmail(subject, body string) {
	// 主题
	m.SetHeader("Subject", subject)

	// 正文
	m.SetBody("text/html", body)

	d := gomail.NewDialer(*serverHost, *serverPort, *fromEmail, *fromEmailPasswd)
	// 发送
	err := d.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}

func Send(subject string, body string, toers []string) {

	myToers := strings.Join(toers, ",")
	myCCers := ""
	// 结构体赋值
	myEmail := &EmailParam{
		ServerHost: *serverHost,
		ServerPort: *serverPort,
		FromEmail:  *fromEmail,
		FromPasswd: *fromEmailPasswd,
		Toers:      myToers,
		CCers:      myCCers,
	}
	initEmail(myEmail)
	sendEmail(subject, body)
}
