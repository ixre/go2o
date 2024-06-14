package tencent

import "github.com/ixre/go2o/core/infrastructure/util/sms"

var _ sms.ISmsProvider = new(TencentSms)

type TencentSms struct {
}

// Name implements sms.ISmsProvider.
func (t *TencentSms) Name() string {
	panic("unimplemented")
}

// Send implements sms.ISmsProvider.
func (t *TencentSms) Send(templateId string, content string, args ...string) error {
	panic("unimplemented")
}
