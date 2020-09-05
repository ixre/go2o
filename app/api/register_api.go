package api

import (
	"errors"
	"fmt"
	"github.com/ixre/gof"
	"github.com/ixre/gof/api"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/registry"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/thrift"
	"go2o/core/service/thrift/auto_gen/rpc/member_service"
	"log"
	"strconv"
	"strings"
	"time"
)

var _ api.Handler = new(RegisterApi)

// 注册接口API
type RegisterApi struct {
	utils
	st storage.Interface
}

func NewRegisterApi() api.Handler {
	st := gof.CurrentApp.Storage()
	return &RegisterApi{
		st: st,
	}
}

func (m RegisterApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"get_token": m.getToken,
		"send_code": m.sendRegisterCode,
		"submit":    m.submit,
	})
}

/**
 * @api {post} /register/submit 用户注册
 * @apiName submit
 * @apiGroup register
 * @apiParam {String} user 用户名
 * @apiParam {String} pwd 密码
 * @apiParam {String} phone 手机号
 * @apiParam {String} token 注册令牌
 * @apiParam {String} check_code 验证码, 如果手机注册时,需要填写
 * @apiParam {String} reg_from 注册来源
 * @apiParam {String} invite_code 邀请码
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (m RegisterApi) submit(ctx api.Context) interface{} {
	user := strings.TrimSpace(ctx.Form().GetString("user"))
	pwd := strings.TrimSpace(ctx.Form().GetString("pwd"))
	phone := strings.TrimSpace(ctx.Form().GetString("phone"))
	regFrom := strings.TrimSpace(ctx.Form().GetString("reg_from"))       // 注册来源
	checkCode := strings.TrimSpace(ctx.Form().GetString("check_code"))   // 验证码
	inviteCode := strings.TrimSpace(ctx.Form().GetString("invite_code")) // 邀请码
	regIp := strings.TrimSpace(ctx.Form().GetString("$user_addr"))       // IP地址
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	if len(token) == 0 || !m.checkRegToken(token) {
		return api.ResponseWithCode(6, "非法注册请求")
	}
	if b := m.compareCheckCode(token, phone, checkCode); !b {
		return api.ResponseWithCode(7, "注册校验码不正确")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		mp := map[string]string{
			"reg_ip":      regIp,
			"reg_from":    regFrom,
			"invite_code": inviteCode,
		}
		r, _ := cli.RegisterMemberV2(thrift.Context, user, pwd, 0, "", phone, "", "", mp)
		if r.ErrCode == 0 {
			//todo: 未生效
			m.signCheckTokenExpires(token)
		}
		return r
	}
	return m.SResult(err)
}

/**
 * @api {post} /register/get_token 获取注册Token
 * @apiName get_token
 * @apiGroup register
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (m RegisterApi) getToken(ctx api.Context) interface{} {
	rd := util.RandString(10)
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time", rd)
	m.st.SetExpire(key, 0, 600)
	return rd
}

// 获取验证码的间隔时间
func (m RegisterApi) getDurationSecond() int64 {
	trans, cli, err := thrift.RegistryServeClient()
	if err == nil {
		val, _ := cli.GetRegistry(thrift.Context, registry.SmsSendDuration)
		trans.Close()
		i, err := strconv.Atoi(val)
		if err != nil {
			log.Println("[ Go2o][ Registry]: parse value error:", err.Error())
		}
		return int64(i)
	}
	return 120
}

// 检查短信验证码是否频繁发送
func (m RegisterApi) checkCodeDuration(token, phone string) error {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time", token)
	nowUnix := time.Now().Unix()
	unix, err := m.st.GetInt64(key)

	log.Println("---", nowUnix, unix, key)

	if err == nil {
		if nowUnix-unix < m.getDurationSecond() {
			return errors.New("请勿在短时间内获取短信验证码!")
		}
	}
	return nil
}

// 标记验证码发送时间
func (m RegisterApi) signCheckCodeSendOk(token string) {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time", token)
	unix := time.Now().Unix()
	log.Println("----save code:", unix)
	m.st.SetExpire(key, unix, 600)
}

// 验证注册令牌是否正确
func (m RegisterApi) checkRegToken(token string) bool {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time", token)
	return m.st.Exists(key)
}

// 将注册令牌标记为过期
func (m RegisterApi) signCheckTokenExpires(token string) {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time", token)
	m.st.Del(key)
}

// 存储校验数据
func (m RegisterApi) saveCheckCodeData(token string, phone string, code string) {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:reg_check_code", token)
	key1 := fmt.Sprintf("sys:go2o:reg:token:%s:reg_check_phone", token)
	m.st.SetExpire(key, code, 600)
	m.st.SetExpire(key1, phone, 600)
}

// 获取校验结果
func (m RegisterApi) compareCheckCode(token, phone string, code string) bool {
	if len(phone) > 0 {
		key := fmt.Sprintf("sys:go2o:reg:token:%s:reg_check_code", token)
		key1 := fmt.Sprintf("sys:go2o:reg:token:%s:reg_check_phone", token)
		ckCode, _ := m.st.GetString(key)
		ckPhone, _ := m.st.GetString(key1)
		if ckPhone == "" || ckPhone != phone {
			return false
		}
		if ckCode == "" || ckCode != code {
			return false
		}
	}
	return true
}

/**
 * @api {post} /register/send_code 发送注册验证码
 * @apiName send_code
 * @apiGroup register
 * @apiParam {String} phone 手机号码
 * @apiParam {String} token 注册令牌
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (m RegisterApi) sendRegisterCode(ctx api.Context) interface{} {
	trans, cli, _ := thrift.RegistryServeClient()
	keys := []string{
		registry.MemberRegisterMustBindPhone,
		registry.SmsRegisterTemplateId,
		registry.EnableDebugMode,
	}
	mp, _ := cli.GetRegistries(thrift.Context, keys)
	trans.Close()
	allowPhoneAsUser := mp[keys[0]]
	debugMode := mp[keys[2]] == "true"
	if allowPhoneAsUser != "true" {
		return api.ResponseWithCode(2, "不允许使用手机号注册")
	}
	phone := strings.TrimSpace(ctx.Form().GetString("phone"))
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	if len(token) == 0 {
		return api.ResponseWithCode(6, "非法注册请求")
	}
	err := m.checkCodeDuration(token, phone)
	if err == nil {
		// 检查手机号码是否被其他人使用
		trans, cli, _ := thrift.MemberServeClient()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Phone, phone)
		trans.Close()
		if memberId <= 0 {
			code := domain.NewCheckCode()
			m.saveCheckCodeData(token, phone, code)
			expiresMinutes := 10
			// 创建参数
			data := map[string]string{
				"code":       code,
				"operation":  "注册会员",
				"minutes":    strconv.Itoa(expiresMinutes),
				"templateId": mp[keys[1]],
			}
			// 构造并发送短信
			trans, cli, _ := thrift.MessageServeClient()
			defer trans.Close()
			n, _ := cli.GetNotifyItem(thrift.Context, "验证手机")
			// 测试环境不发送短信
			if debugMode {
				return api.ResponseWithCode(3, "【测试】短信验证码为:"+code)
			}
			// 发送短信
			r, _ := cli.SendPhoneMessage(thrift.Context, phone, n.Content, data)
			if r.ErrCode == 0 {
				m.signCheckCodeSendOk(code) // 标记为已发送
			} else {
				log.Println("[ Go2o][ Sms]: 验证码发送失败:", r.ErrMsg)
				return api.ResponseWithCode(3, "验证码发送失败")
			}
		} else {
			return api.ResponseWithCode(1, "手机号码已注册")
		}
	}
	return api.NewResponse(nil)
}
