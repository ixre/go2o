package sp

import (
	"testing"

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	s "github.com/ixre/go2o/core/infrastructure/util/sms"
	"github.com/ixre/go2o/core/inject"
	_ "github.com/ixre/go2o/tests"
)

func TestSendTencentSms(t *testing.T) {
	repo := inject.GetMessageRepo()
	templates := repo.NotifyRepo().GetAllNotifyTemplate()
	temp := collections.FindArray(templates, func(t *mss.NotifyTemplate) bool {
		return t.TplType == 2 && t.TplCode == mss.SMS_CHECK_CODE
	})
	err := s.Send(s.Template{
		ProviderCode:    temp.SpCode,
		TemplateContent: "",
		TemplateId:      temp.SpTid,
	}, "13068686358", "9996", "5")
	if err != nil {
		t.Error(err)
	}
}
