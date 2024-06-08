package domain

import (
	"strconv"
	"testing"
	"time"

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/message/notify"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/inject"
)

// 测试配置短信服务商参数
func TestConfigureSms(t *testing.T) {
	provider := mss.CHUANGLAN
	manager := inject.GetMessageRepo().NotifyManager()
	err := manager.SaveSmsApiPerm(&notify.SmsApiPerm{
		Provider:   int(provider),
		Key:        "N42622266620",
		Secret:     "X34Mvw5f5db",
		Signature:  "【go2o】",
		TemplateId: "",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	repo := inject.GetRegistryRepo()
	ir := repo.Get(registry.SmsDefaultProvider)
	err = ir.Update(strconv.Itoa(int(provider)))
	if err == nil {
		err = ir.Save()
	}
	if err != nil {
		t.Error(err)
	}
	v := repo.Get(registry.SmsDefaultProvider).IntValue()
	if v != int(provider) {
		t.Log("未保存成功")
		t.Fail()
	}
}

// 测试发送短信
func TestSendPhoneMessage(t *testing.T) {
	templatId := notify.SMS_CHECK_CODE
	manager := inject.GetMessageRepo().NotifyManager()

	err := manager.SendPhoneMessage("13162222872",
		notify.PhoneMessage("测试短信:你本次进行{action}的验证码为{1},有效期为:{minites}"),
		[]string{"注册账户", "3101", "30"},
		templatId,
	)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(3000)
}
