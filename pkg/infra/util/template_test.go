package util

import (
	"testing"
)

func TestResolveMessageVariables(t *testing.T) {
	//msg := `Hello ${Name}, your age is ${Age}, Thank you ${Name}`
	msg := `您正在注册成为商户，请点击以下链接完成注册： <br /><a href="${注册链接}">${注册链接}</a>
		 			<br />
		 		此链接有效期为${有效时间}分钟`
	s := ResolveMessage(msg, []string{"Tommy", "20"})
	t.Log(s)
}
func TestResolveMessageNumVariables(t *testing.T) {
	msg := `Hello {0}, your age is {1}, Thank you {0}`
	s := ResolveMessage(msg, []string{"Tommy", "20"})
	t.Log(s)
}
