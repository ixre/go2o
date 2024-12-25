package sp

import (
	"testing"

	"github.com/ixre/go2o/core/infrastructure/util/smtp"
)

// TestSendMerchantCheckCodeMail 是一个测试函数，用于测试发送商户验证码邮件的功能
func TestSendMerchantCheckCodeMail(t *testing.T) {
	err := smtp.SendMail("[重要]商户注册链接", []string{
		"959398298@qq.com",
	}, `
	
		<div style="font-size:16px;color:#333">测试邮件</div>
	`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
