package helper

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

var SendEmail string
var MailPassword string

func init() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("fail to read config file" + err.Error())
	}

	SendEmail = viper.GetString("send_email")
	MailPassword = viper.GetString("mail_password")
}

// 发送验证码
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Get <" + SendEmail + ">"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送"
	e.HTML = []byte("您的验证码:<b>" + code + "</b>")
	return e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", SendEmail, MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}

func GetCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	res := ""
	for i := 0; i < 6; i++ {
		res += strconv.Itoa(rand.Intn(10))
	}
	return res
}
