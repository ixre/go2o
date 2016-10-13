/**
 * Copyright 2015 @ z3q.net.
 * name : mss_test
 * author : jarryliu
 * date : 2016-07-06 20:22
 * description :
 * history :
 */
package testing

import (
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/infrastructure/tool/sms"
	"go2o/core/repository"
	"go2o/core/testing/include"
	"testing"
)

func TestMssSendSms(t *testing.T) {
	app := include.GetApp()
	db := app.Db()
	sto := app.Storage()
	nRep := repository.NewNotifyRep(db)
	vRep := repository.NewValueRep(db, sto)
	rep := repository.NewMssRep(db, nRep, vRep)

	data := map[string]interface{}{}
	data = sms.AppendCheckPhoneParams(1, data)
	err := rep.NotifyManager().SendPhoneMessage("18616999822",
		notify.PhoneMessage("您 好啊"), data)
	if err != nil {
		t.Fatal(err)
	}
}
