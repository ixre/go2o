/**
 * Copyright 2015 @ z3q.net.
 * name : mail.go
 * author : jarryliu
 * date : 2015-07-26 20:14
 * description :
 * history :
 */
package mss

import (
	"encoding/base64"
	"fmt"
	"github.com/ixre/gof"
	"go2o/core/variable"
	"net/smtp"
	"strings"
)

var (
	loaded               bool   = false
	EMAIL_SERVER         string = ""
	EMAIL_HOST           string = ""
	EMAIL_CREDENTIAL_USR        = ""
	EMAIL_CREDENTIAL_PWD        = ""
	EMAIL_FROM                  = ""
)

func SendMail(server, host, usr, pwd, from string, subject string, to []string, body []byte) error {
	auth := smtp.PlainAuth("", usr, pwd, host)
	header := make(map[string]string)
	header["From"] = from
	header["To"] = strings.Join(to, ";")
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=utf-8"
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(body)

	return smtp.SendMail(server, auth, from, to, []byte(message))
}

// 使用默认的配置发送邮件
func SendMailWithDefaultConfig(subject string, to []string, body []byte) error {
	if !loaded {
		cfg := gof.CurrentApp.Config()
		EMAIL_HOST = cfg.GetString(variable.SmtpHost)
		EMAIL_SERVER = fmt.Sprintf("%s:%d", EMAIL_HOST, cfg.GetInt(variable.SmtpPort))
		EMAIL_CREDENTIAL_USR = cfg.GetString(variable.SmtpCreUser)
		EMAIL_CREDENTIAL_PWD = cfg.GetString(variable.SmtpCrePwd)
		EMAIL_FROM = cfg.GetString(variable.SmtpFrom)
		loaded = true
	}
	return SendMail(EMAIL_SERVER, EMAIL_HOST, EMAIL_CREDENTIAL_USR, EMAIL_CREDENTIAL_PWD,
		EMAIL_FROM, subject, to, body)
}
