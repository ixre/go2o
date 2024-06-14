package sms

// 短信服务商
var smsProviders []ISmsProvider

type ISmsProvider struct {
}

type Template struct {
	// 服务商代码
	ProviderCode string
	// 模板内容,如果传递了TemplateId,使用短信服务商申请的模板内容发送
	TemplateContent string
	// 模板ID,如果自定义短信，则不用传ID
	TemplateId string
}

func RegisterProvider(p ISmsProvider) {

}

// 发送短信
func Send(t Template, phoneNum string, params []string) error {
	//c := ResolveMessage(content, params)

	// todo: 旧的短信需要重新实现
	// if setting.Signature != "" && !strings.Contains(content, setting.Signature) {
	// 	content = setting.Signature + content
	// }
	// c := ResolveMessage(content, params)
	// switch mss.SmsProvider(setting.Provider) {
	// case mss.HTTP:
	// 	return sendPhoneMsgByHttpApi(setting, phoneNum, c, params, "")
	// case mss.TECENT_CLOUD:
	// 	return errors.New("not implemented")
	// case mss.ALIYUN:
	// 	templateName := ""
	// 	return aliyu.SendSms(setting.Key,
	// 		setting.Secret, phoneNum,
	// 		content, params,
	// 		templateName, "")
	// case mss.CHUANGLAN:
	// 	return cl253.SendMsgToMobile(setting.Key, setting.Secret, phoneNum, c)
	// }
	// return nil
	return nil
}
