package smtp

import (
	"errors"
	"log"

	"github.com/ixre/go2o/core/infrastructure/logger"
	"gopkg.in/gomail.v2"
)

type SmtpConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

var _defaultCfg *SmtpConfig

// func sendMail(subject string, to []string, body []byte, cfg *SmtpConfig) error {
// 	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
// 	header := make(map[string]string)
// 	header["To"] = strings.Join(to, ";")
// 	header["Subject"] = fmt.Sprintf("= ?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
// 	header["MIME-Version"] = "1.0"
// 	header["Content-Type"] = "text/html; charset=utf-8"
// 	header["Content-Transfer-Charset"] = "base64"

// 	message := ""
// 	for k, v := range header {
// 		message += fmt.Sprintf("%s: %s\r\n", k, v)
// 	}
// 	message += "\r\n" + base64.StdEncoding.EncodeToString(body)

// 	return smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
// 		auth, cfg.From, to, []byte(message))
// }

// 发送邮件
func sendMailWithGoMail(subject string, to []string, body string, cfg *SmtpConfig) error {
	// 设置SMTP服务器信息
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)
	// 创建一个新的邮件消息
	m := gomail.NewMessage()
	if len(cfg.From) > 0 {
		// 使用别名
		m.SetHeader("From", m.FormatAddress(cfg.User, cfg.From))
	} else {
		m.SetHeader("From", cfg.User)
	}
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	// 设置邮件正文为HTML
	m.SetBody("text/html", body)
	log.Println("---", body)
	// 发送邮件
	return d.DialAndSend(m)
}

// 发送邮件
func SendMail(subject string, to []string, body string) error {
	if _defaultCfg == nil {
		return errors.New("未配置邮箱服务器")
	}
	return sendMailWithGoMail(subject, to, body, _defaultCfg)
}

func Configure(cfg *SmtpConfig) {
	if cfg.Port == 0 || cfg.Host == "" || cfg.User == "" || cfg.Password == "" {
		logger.Warn("邮箱服务器设置不正确,将影响邮件发送功能")
		return
	}
	_defaultCfg = cfg
}
