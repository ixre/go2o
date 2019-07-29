package api

import (
	"errors"
	"fmt"
	"github.com/ixre/goex/echox"
	"github.com/ixre/gof"
	"github.com/ixre/gof/api"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/registry"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/thrift"
	"log"
	"strconv"
	"strings"
	"time"
)

var _ api.Handler = new(RegisterApi)

type RegisterApi struct {
	apiUtil
	st storage.Interface
}

func NewRegisterApi()api.Handler{
	st := gof.CurrentApp.Storage()
	return &RegisterApi{
		st:      st,
	}
}

func (m RegisterApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"get_token":m.getToken,
		"submit":m.submit,
	})
}



/**
 * @api {post} /register/submit 用户注册
 * @apiName submit
 * @apiGroup register
 * @apiParam {String} user 用户名
 * @apiParam {String} pwd 密码
 * @apiParam {String} phone 手机号
 * @apiParam {String} reg_from 注册来源
 * @apiParam {String} invite_code 邀请码
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (m RegisterApi) submit(ctx api.Context) interface{} {
	user := ctx.Form().GetString("user")
	pwd := ctx.Form().GetString("pwd")
	phone := ctx.Form().GetString("phone")
	regFrom := ctx.Form().GetString("reg_from")       // 注册来源
	inviteCode := ctx.Form().GetString("invite_code") // 邀请码
	regIp := ctx.Form().GetString("$user_ip_addr")    // IP地址
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		mp := map[string]string{
			"reg_ip":      regIp,
			"reg_from":    regFrom,
			"invite_code": inviteCode,
		}
		r, _ := cli.RegisterMemberV2(thrift.Context, user, pwd, 0, "", phone, "", "", mp)
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
func (m RegisterApi) getToken(ctx api.Context)interface{}{
	rd := util.RandString(10)
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time",rd)
	m.st.SetExpire(key,rd,600)
	return rd
}


// 获取验证码的间隔时间
func (m RegisterApi) getDuractionSecond() int64 {
	trans, cli, err := thrift.FoundationServeClient()
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
func (m RegisterApi) checkCodeDuration_Reg(token,phone string) error {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time",token)
	nowUnix := time.Now().Unix()
	unix,err := m.st.GetInt64(key)
	if err == nil {
		if nowUnix - unix < m.getDuractionSecond(){
			return errors.New("请勿在短时间内获取短信验证码!")
		}
	}
	return nil
}



// 存储校验数据
func (m RegisterApi) saveCheckData_Reg(token string,phone string, code string) {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:reg_check_code", token)
	key1 := fmt.Sprintf("sys:go2o:reg:token:%s:reg_check_phone", token)
	m.st.SetExpire(key, code, 600)
	m.st.SetExpire(key1, phone, 600)
}

// 标记验证码发送时间
func (m RegisterApi) signCheckCodeSend_Reg(c *echox.Context) {
	ss := c.Session
	ss.Set("reg_last_unix", time.Now().Unix()) //最后的发送时间
	ss.Save()
}

// 发送验证码
func (m RegisterApi) SendRegisterCode(ctx api.Context) interface{} {
	trans, cli, _ := thrift.FoundationServeClient()
	keys := []string{
		registry.MemberRegisterMustBindPhone,
		registry.SmsRegisterTemplateId,
		registry.EnableDebugMode,
	}
	mp, _ := cli.GetRegistries(thrift.Context, keys)
	trans.Close()
	allowPhoneAsUser := mp[keys[0]]
	debugMode := mp[keys[2]] == "true"
	if  allowPhoneAsUser!= "true" {
		return api.ResponseWithCode(2,"不允许使用手机号注册")
	}
	phone := strings.TrimSpace(ctx.Form().GetString("phone"))
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	if len(token) == 0{
		return api.ResponseWithCode(6,"非法注册请求")
	}
	err := m.checkCodeDuration_Reg(token,phone)
	if err == nil {
		// 检查手机号码是否被其他人使用
		trans, cli, _ := thrift.MemberServeClient()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Phone, phone)
		trans.Close()
		if memberId <= 0 {
			code := domain.NewCheckCode()
			m.saveCheckData_Reg(token,phone,code)
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
				return api.ResponseWithCode(3,"【测试】短信验证码为:" + code)
			}
			// 发送短信
			r, _ := cli.SendPhoneMessage(thrift.Context, phone, n.Content, data)
			if r.ErrCode == 0 {
				m.signCheckCodeSend_Reg(code) // 标记为已发送
			} else {
				log.Println("[ Go2o][ Sms]: 验证码发送失败:", r.ErrMsg)
				return api.ResponseWithCode(3,"验证码发送失败")
			}
		} else {
			return api.ResponseWithCode(1,"手机号码已注册")
		}
	}
	return api.NewResponse(map[string]string{})
}

