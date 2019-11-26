package api

import (
	"errors"
	"fmt"
	"github.com/ixre/gof"
	"github.com/ixre/gof/api"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/registry"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/auto_gen/rpc/message_service"
	"go2o/core/service/thrift"
	"log"
	"strconv"
	"strings"
	"time"
)

var _ api.Handler = new(PassportApi)

var (
	operationArr = []string{"找回密码", "重置密码", "绑定手机"}
	//2分钟后才可重发验证码
	timeOutUnix int64 = 120 //等于:time.Unix(120,0).Unix()
)

type PassportApi struct {
utils
	st storage.Interface
}

func NewPassportApi() api.Handler {
	st := gof.CurrentApp.Storage()
	return &PassportApi{
		st: st,
	}
}

func (m PassportApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"get_token":       m.getToken,
		"send_code":       m.sendCode,
		"compare_code":    m.compareCode,
		"modify_pwd":      m.modifyPwd,
		"reset_pwd":       m.resetPwd,
		"trade_pwd":       m.tradePwd,
		"reset_trade_pwd": m.resetTradePwd,
	})
}

// 根据输入的凭据获取会员编号
func (m PassportApi) checkMemberBasis(ctx api.Context) (string, member_service.ECredentials, error) {
	acc := strings.TrimSpace(ctx.Form().GetString("account"))          //账号、手机或邮箱
	credTypeId, err := strconv.Atoi(ctx.Form().GetString("cred_type")) //账号类型
	if err != nil {
		return acc, member_service.ECredentials_User, err
	}
	credType := member_service.ECredentials(credTypeId)
	if len(acc) == 0 {
		return acc, credType, errors.New("信息不完整")
	}
	return acc, credType, nil
}

// 根据发送的校验码类型获取用户凭据类型
func (h PassportApi) parseMessageChannel(credType member_service.ECredentials) message_service.EMessageChannel {
	if credType == member_service.ECredentials_Email {
		return message_service.EMessageChannel_EmailMessage
	}
	return message_service.EMessageChannel_SmsMessage
}

// 标记验证码发送时间
func (h PassportApi) signCodeSendInfo(token string) {
	prefix := "sys:go2o:pwd:token"
	// 最后的发送时间
	unix := time.Now().Unix()
	h.st.SetExpire(fmt.Sprintf("%s:%s:last-time", prefix, token), unix, 600)
	// 验证码校验成功
	h.st.SetExpire(fmt.Sprintf("%s:%s:check_ok", prefix, token), 0, 600)
	// 清除记录的会员编号
	h.st.Del(fmt.Sprintf("%s:%s:member_id", prefix, token))
}

// 获取校验结果
func (p PassportApi) GetCodeVerifyResult(token string) (int64, bool) {
	prefix := "sys:go2o:pwd:token"
	checkKey := fmt.Sprintf("%s:%s:check_ok", prefix, token)
	v, err := p.st.GetInt64(checkKey)
	//验证码校验成功
	if err == nil && v == 1 {
		mmKey := fmt.Sprintf("%s:%s:member_id", prefix, token)
		memberId, err := p.st.GetInt64(mmKey)
		return memberId, err == nil
	}
	return 0, false
}

// 清理验证码校验结果
func (p PassportApi) resetCodeVerifyResult(token string) {
	prefix := "sys:go2o:pwd:token"
	checkKey := fmt.Sprintf("%s:%s:check_ok", prefix, token)
	p.st.Del(checkKey)
}

// 设置校验成功
func (p PassportApi) setCodeVerifySuccess(token string, memberId int64) {
	prefix := "sys:go2o:pwd:token"
	checkKey := fmt.Sprintf("%s:%s:check_ok", prefix, token)
	mmKey := fmt.Sprintf("%s:%s:member_id", prefix, token)
	p.st.SetExpire(checkKey, 1, 600)     // 验证码校验成功
	p.st.SetExpire(mmKey, memberId, 600) // 记录会员编号
}

// 验证令牌是否正确
func (m PassportApi) checkToken(token string) bool {
	key := fmt.Sprintf("sys:go2o:pwd:token:%s:last-time", token)
	return m.st.Exists(key)
}

/**
 * @api {post} /passport/get_token 获取注册Token
 * @apiName get_token
 * @apiGroup passport
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (h PassportApi) getToken(ctx api.Context) interface{} {
	rd := util.RandString(10)
	key := fmt.Sprintf("sys:go2o:pwd:token:%s:last-time", rd)
	h.st.SetExpire(key, 0, 600)
	return rd
}

/**
 * @api {post} /passport/send_code 发送验证码
 * @apiName send_code
 * @apiGroup passport
 * @apiParam {String} account 账号
 * @apiParam {Int} cred_type 账号类型,1:用户名 3:邮件 4:手机号
 * @apiParam {String} token 令牌
 * @apiParam {Int} op 验证码场景:0:找回密码, 1:重置密码 2:绑定手机
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (h PassportApi) sendCode(ctx api.Context) interface{} {
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	if len(token) == 0 || !h.checkToken(token) {
		return api.ResponseWithCode(6, "令牌无效")
	}
	operation, _ := strconv.Atoi(ctx.Form().GetString("op")) //操作
	account, credType, err := h.checkMemberBasis(ctx)
	if err != nil {
		return api.ResponseWithCode(2, err.Error())
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, credType, account)
		if memberId <= 0 {
			return api.ResponseWithCode(1, member.ErrNoSuchMember.Error())
		}
		err = h.checkCodeDuration(token, account)
		if err == nil {
			var msgChan = h.parseMessageChannel(credType)
			r, _ := cli.SendCode(thrift.Context, memberId,
				operationArr[operation], msgChan)
			code := r.Data["code"]
			if r.ErrCode == 0 {
				h.signCodeSendInfo(token) // 标记为已发送
			} else {
				log.Println("[ Go2o][ Error]: 发送会员验证码失败:", r.ErrMsg)
				err = errors.New("发送验证码失败")
			}
			keys := []string{
				registry.EnableDebugMode,
			}
			trans, cli, _ := thrift.RegistryServeClient()
			mp, _ := cli.GetRegistries(thrift.Context, keys)
			trans.Close()
			debugMode := mp[keys[0]] == "true"
			if debugMode && len(code) != 0 {
				return api.ResponseWithCode(3, "【测试】短信验证码为:"+code)
			}
		}
	}
	if err != nil {
		return api.ResponseWithCode(1, err.Error())
	}
	return api.NewResponse(map[string]string{})
}

/**
 * @api {post} /passport/compare_code 校验验证码
 * @apiName compare_code
 * @apiGroup passport
 * @apiParam {String} account 账号
 * @apiParam {Int} cred_type 账号类型,1:站内信 3:邮箱 4:短信
 * @apiParam {String} token 令牌
 * @apiParam {String} check_code 校验码
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (h PassportApi) compareCode(ctx api.Context) interface{} {
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	if len(token) == 0 {
		return api.ResponseWithCode(6, "非法注册请求")
	}
	account, credType, err := h.checkMemberBasis(ctx)
	if err != nil {
		return api.ResponseWithCode(2, err.Error())
	}
	code := ctx.Form().GetString("check_code")
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, credType, account)
		r, _ := cli.CompareCode(thrift.Context, memberId, code)
		if r.ErrCode == 0 {
			h.setCodeVerifySuccess(token, memberId)
		} else {
			err = errors.New(r.ErrMsg)
		}
	}
	if err != nil {
		return api.ResponseWithCode(1, err.Error())
	}
	return api.NewResponse(map[string]string{})
}

/**
 * @api {post} /passport/reset_pwd 重置密码
 * @apiName reset_pwd
 * @apiGroup passport
 * @apiParam {String} account 账号
 * @apiParam {Int} cred_type 账号类型,1:用户名 3:邮件 4:手机号
 * @apiParam {String} token 令牌
 * @apiParam {String} pwd md5编码后的密码
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (h PassportApi) resetPwd(ctx api.Context) interface{} {
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	pwd := strings.TrimSpace(ctx.Form().GetString("pwd"))
	if len(token) == 0 || !h.checkToken(token) {
		return api.ResponseWithCode(6, "请求已过期")
	}
	account, credType, err := h.checkMemberBasis(ctx)
	if err != nil {
		return api.ResponseWithCode(2, err.Error())
	}
	// 验证校验码是否正确
	memberId, b := h.GetCodeVerifyResult(token)
	if !b {
		return api.ResponseWithCode(2, "验证无效")
	}
	// 验证会员是否匹配
	if err := h.checkMemberMatch(account, credType, memberId); err != nil {
		return api.ResponseWithCode(1, err.Error())
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.ModifyPwd(thrift.Context, memberId, "", pwd)
		if r.ErrCode != 0 {
			return api.ResponseWithCode(int(r.ErrCode), r.ErrMsg)
		}
	}
	h.resetCodeVerifyResult(token)
	return api.NewResponse(map[string]string{})
}

/**
 * @api {post} /passport/modify_pwd 修改密码
 * @apiName modify_pwd
 * @apiGroup passport
 * @apiParam {String} account 账号
 * @apiParam {Int} cred_type 账号类型,1:用户名 3:邮件 4:手机号
 * @apiParam {String} token 令牌
 * @apiParam {String} pwd md5编码后的密码
 * @apiParam {String} old_pwd md5编码后的旧密码
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 *
 */
func (h PassportApi) modifyPwd(ctx api.Context) interface{} {
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	pwd := strings.TrimSpace(ctx.Form().GetString("pwd"))
	oldPwd := strings.TrimSpace(ctx.Form().GetString("old_pwd"))
	if len(token) == 0 || !h.checkToken(token) {
		return api.ResponseWithCode(6, "请求已过期")
	}
	account, credType, err := h.checkMemberBasis(ctx)
	if err != nil {
		return api.ResponseWithCode(2, err.Error())
	}
	// 验证校验码是否正确
	memberId, b := h.GetCodeVerifyResult(token)
	if !b {
		return api.ResponseWithCode(2, "验证无效")
	}
	// 验证会员是否匹配
	if err := h.checkMemberMatch(account, credType, memberId); err != nil {
		return api.ResponseWithCode(1, err.Error())
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.ModifyPwd(thrift.Context, memberId, oldPwd, pwd)
		if r.ErrCode != 0 {
			return api.ResponseWithCode(int(r.ErrCode), r.ErrMsg)
		}
	}
	h.resetCodeVerifyResult(token)
	return api.NewResponse(map[string]string{})
}

/**
 * @api {post} /passport/trade_pwd 修改交易密码
 * @apiName trade_pwd
 * @apiGroup passport
 * @apiParam {String} account 账号
 * @apiParam {Int} cred_type 账号类型,1:用户名 3:邮件 4:手机号
 * @apiParam {String} token 令牌
 * @apiParam {String} pwd md5编码后的密码
 * @apiParam {String} old_pwd md5编码后的旧密码,如果没有设置交易密码,则为空
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 *
 */
func (h PassportApi) tradePwd(ctx api.Context) interface{} {
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	pwd := strings.TrimSpace(ctx.Form().GetString("pwd"))
	oldPwd := strings.TrimSpace(ctx.Form().GetString("old_pwd"))
	if len(token) == 0 || !h.checkToken(token) {
		return api.ResponseWithCode(6, "请求已过期")
	}
	account, credType, err := h.checkMemberBasis(ctx)
	if err != nil {
		return api.ResponseWithCode(2, err.Error())
	}
	// 验证校验码是否正确
	memberId, b := h.GetCodeVerifyResult(token)
	if !b {
		return api.ResponseWithCode(2, "验证无效")
	}
	// 验证会员是否匹配
	if err := h.checkMemberMatch(account, credType, memberId); err != nil {
		return api.ResponseWithCode(1, err.Error())
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.ModifyTradePwd(thrift.Context, memberId, oldPwd, pwd)
		if r.ErrCode != 0 {
			return api.ResponseWithCode(int(r.ErrCode), r.ErrMsg)
		}
	}
	h.resetCodeVerifyResult(token)
	return api.NewResponse(map[string]string{})
}

/**
 * @api {post} /passport/reset_trade_pwd 重置交易密码
 * @apiName reset_trade_pwd
 * @apiGroup passport
 * @apiParam {String} account 账号
 * @apiParam {Int} cred_type 账号类型,1:用户名 3:邮件 4:手机号
 * @apiParam {String} token 令牌
 * @apiParam {String} pwd md5编码后的密码
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 *
 */
func (h PassportApi) resetTradePwd(ctx api.Context) interface{} {
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	pwd := strings.TrimSpace(ctx.Form().GetString("pwd"))
	if len(token) == 0 || !h.checkToken(token) {
		return api.ResponseWithCode(6, "请求已过期")
	}
	account, credType, err := h.checkMemberBasis(ctx)
	if err != nil {
		return api.ResponseWithCode(2, err.Error())
	}
	// 验证校验码是否正确
	memberId, b := h.GetCodeVerifyResult(token)
	if !b {
		return api.ResponseWithCode(2, "验证无效")
	}
	// 验证会员是否匹配
	if err := h.checkMemberMatch(account, credType, memberId); err != nil {
		return api.ResponseWithCode(1, err.Error())
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.ModifyTradePwd(thrift.Context, memberId, "", pwd)
		if r.ErrCode != 0 {
			return api.ResponseWithCode(int(r.ErrCode), r.ErrMsg)
		}
	}
	h.resetCodeVerifyResult(token)
	return api.NewResponse(map[string]string{})
}

// 获取验证码的间隔时间
func (m PassportApi) getDurationSecond() int64 {
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
func (m PassportApi) checkCodeDuration(token, phone string) error {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time", token)
	nowUnix := time.Now().Unix()
	unix, err := m.st.GetInt64(key)
	if err == nil {
		if nowUnix-unix < m.getDurationSecond() {
			return errors.New("请勿在短时间内获取短信验证码!")
		}
	}
	return nil
}

// 标记验证码发送时间
func (m PassportApi) signCheckCodeSendOk(token string) {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time", token)
	unix := time.Now().Unix()
	log.Println("----save code:", unix)
	m.st.SetExpire(key, unix, 600)
}

// 验证注册令牌是否正确
func (m PassportApi) checkRegToken(token string) bool {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time", token)
	_, err := m.st.GetInt64(key)
	return err == nil
}

// 将注册令牌标记为过期
func (m PassportApi) signCheckTokenExpires(token string) {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:last-time", token)
	m.st.Del(key)
}

// 存储校验数据
func (m PassportApi) saveCheckCodeData(token string, phone string, code string) {
	key := fmt.Sprintf("sys:go2o:reg:token:%s:reg_check_code", token)
	key1 := fmt.Sprintf("sys:go2o:reg:token:%s:reg_check_phone", token)
	m.st.SetExpire(key, code, 600)
	m.st.SetExpire(key1, phone, 600)
}

// 验证会员是否匹配
func (m PassportApi) checkMemberMatch(account string, credType member_service.ECredentials, memberId int64) error {
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		mid, _ := cli.SwapMemberId(thrift.Context, credType, account)
		if mid <= 0 {
			return member.ErrNoSuchMember
		}
		if mid != memberId {
			return errors.New("member not match")
		}
	}
	return err
}
