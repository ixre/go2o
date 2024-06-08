/**
 * Copyright 2015 @ 56x.net.
 * name : mss_test
 * author : jarryliu
 * date : 2016-07-06 20:22
 * description :
 * history :
 */
package domain

import (
	"testing"

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"

	"github.com/ixre/go2o/core/inject"
)

func TestMssSendSms(t *testing.T) {
	nRepo := inject.GetMessageRepo()
	registryRepo := inject.GetRegistryRepo()
	nm := nRepo.NotifyManager()
	re := registryRepo.Get(registry.SmsDefaultProvider)
	re.Update("http")
	re.Save()
	re = registryRepo.Get(registry.SmsRegisterTemplateId)
	re.Update("8332")
	re.Save()
	err := nm.SaveSmsApiPerm(&mss.SmsApiPerm{
		Extra: &mss.SmsExtraSetting{
			ApiUrl: "https://api.zhuanxinyun.com/api/v2/sendSms.json",
			Params: "appKey={key}&appSecret={secret}&phones={phone}&content={content}" +
				"&batchNum={stamp}&templateId={templateId}",
			Method:       "POST",
			Charset:      "utf-8",
			SuccessChars: "errorCode\":\"000000\"",
		},
		Key:       "NUV2LeZr4c6Ta2tdMHK1AfSsaut1Jscf",
		Secret:    "9f5946bb1dac95e87ef69d7e5e8e0a08",
		Signature: "",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}

func TestSendSms(t *testing.T) {
	//data := map[string]interface{}{
	//	"templateId":8332,
	//}
	//err = nm.SendPhoneMessage("13162222872",
	//	"您正在重置密码, 本次验证码为: 3366, 有效期5分钟, 请确保是您本人操作!",
	//	data)
	//if err != nil {
	//	t.Fatal(err)
	//}
}
