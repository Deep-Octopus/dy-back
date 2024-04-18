package utils

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"math/rand"
)

var EmailVerificationCode = make(map[string]string, 10)

func SendVerificationCode(targetEmail string, code string) error {
	//from := CONF.Smtp.From
	//password := CONF.Smtp.Password
	//smtpHost := CONF.Smtp.Host
	//smtpPort := CONF.Smtp.Port
	//
	//auth := smtp.PlainAuth("", from, password, smtpHost)
	//
	//msg := fmt.Sprintf("Subject: Email Verification Code\r\n\r\n%s", code)
	//to := []string{email}
	//err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	//return err

	server := mail.NewSMTPClient()
	server.Host = CONF.Smtp.Host
	server.Port = CONF.Smtp.Port
	server.Username = CONF.Smtp.From
	server.Password = CONF.Smtp.Password
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(CONF.Smtp.From)
	email.AddTo(targetEmail)
	//email.AddCc("another_you@example.com")
	email.SetSubject("抖音：")

	htmlStr := "<html>验证码:" + code + "</html>"
	email.SetBody(mail.TextHTML, htmlStr) //发送html信息
	//email.AddAttachment("super_cool_file.png") // 附件

	// Send email
	err = email.Send(smtpClient)
	return err
}

func GenerateVerificationCode() string {
	// 生成随机的6位验证码
	const charset = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
