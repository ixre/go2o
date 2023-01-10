package domain

import (
	"testing"
	"time"

	"github.com/ixre/go2o/core/domain/interface/mss/notify"
	"github.com/ixre/go2o/tests/ti"
)

func TestSendPhoneMessage(t *testing.T) {
	templatId := ""
	manager := ti.Factory.GetMssRepo().NotifyManager()
	err := manager.SendPhoneMessage("13162222872",
		notify.PhoneMessage("测试短信消息"),
		[]string{"注册账户", "jarrysix"},
		templatId,
	)
	if err != nil{
		t.Error(err)
	}
	time.Sleep(3000)
}
