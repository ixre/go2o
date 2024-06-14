package tencent

import (
	"fmt"
	"strconv"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/util/sms"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	ts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

var _ sms.ISmsProvider = new(TencentSms)
var secretIdKey = "sp_tencent_secret_id"
var secretKey = "sp_tencent_secret_key"
var smsAppId = "sp_tencent_sms_appid"
var smsSignName = "sp_tencent_sms_sign_name"

type TencentSms struct {
	repo registry.IRegistryRepo
}

func NewTencentSms(repo registry.IRegistryRepo) *TencentSms {
	initConfig(repo)
	return &TencentSms{
		repo: repo,
	}
}

func initConfig(repo registry.IRegistryRepo) {
	repo.CreateUserKey(secretIdKey, "-", "腾讯云API接口SecretId")
	repo.CreateUserKey(secretKey, "-", "腾讯云API接口SecretKey")
	repo.CreateUserKey(smsAppId, "-", "腾讯云短信APPId")
	repo.CreateUserKey(smsSignName, "[短信签名]", "腾讯云短信签名")
}

// Name implements sms.ISmsProvider.
func (t *TencentSms) Name() string {
	return "TENCENT"
}

// SendContent implements sms.ISmsProvider.
func (t *TencentSms) SendContent(phone string, content string) error {
	return fmt.Errorf("%s not support content sms", t.Name())
}

// Send implements sms.ISmsProvider.
func (t *TencentSms) Send(phone string, templateId string, args ...string) error {
	req := ts.NewSendSmsRequest()
	req.PhoneNumberSet = common.StringPtrs([]string{
		phone,
	})
	req.TemplateId = common.StringPtr(templateId)
	req.TemplateParamSet = common.StringPtrs(args)
	req.TemplateId = common.StringPtr(templateId)
	return t.send(req)
}

func (t *TencentSms) send(req *ts.SendSmsRequest) error {
	// 初始化用户身份信息（secretId, secretKey）
	secretId, _ := t.repo.GetValue(secretIdKey)
	secretKey, _ := t.repo.GetValue(secretKey)
	appId, _ := t.repo.GetValue(smsAppId)
	signName, _ := t.repo.GetValue(smsSignName)
	debug, _ := t.repo.GetValue(registry.EnableDebugMode)

	// 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId, secretKey
	credential := common.NewCredential(secretId, secretKey)
	// 配置签名和应用Id
	req.SmsSdkAppId = common.StringPtr(appId)
	req.SignName = common.StringPtr(signName)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	cpf.SignMethod = "HmacSHA1"

	// 实例化SMS的client对象
	smsClient, err := ts.NewClient(credential, "ap-guangzhou", cpf)
	if err == nil {
		// 通过client对象调用SendSms方法发起请求
		ret, err := smsClient.SendSms(req)
		if err != nil {
			return err
		}
		if b, _ := strconv.ParseBool(debug); b {

			// 输出返回的响应内容: 如果发送未收到,可能是超出每日发送数量
			fmt.Printf("发送腾讯云短信: %+v\n", ret.ToJsonString())
		}
	}
	return err

}
