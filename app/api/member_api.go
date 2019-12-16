package api

import (
	"encoding/json"
	"fmt"
	"github.com/ixre/gof/api"
	"github.com/ixre/gof/types"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/registry"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"strconv"
	"strings"
)

var _ api.Handler = new(MemberApi)

var provider = map[string]string{
	"alipay": "支付宝",
	"wepay":  "微信支付",
	"unipay": "云闪付",
}

type MemberApi struct {
	utils
}

func (m MemberApi) Process(fn string, ctx api.Context) *api.Response {
	var memberId int64
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) > 0{
		memberId, _ = rsi.MemberService.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
	}
	switch fn {
	case "order_summary":return m.orderSummary(ctx,memberId)
	case "orders_quantity":return m.ordersQuantity(ctx,memberId)
	}
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"login":           m.login,
		"get":             m.getMember,
		"account":         m.account,
		"profile":         m.profile,
		"checkToken":      m.checkToken,
		"check_token":     m.checkToken,
		"complex":         m.complex,
		"bankcard":        m.bankcard,
		"invites":         m.invites,
		"receipts_code":   m.receiptsCode,
		"save_receipts":   m.saveReceiptsCode,
		"toggle_receipts": m.toggleReceipts,
	})
}

// 登录
func (m MemberApi) login(ctx api.Context) interface{} {
	form := ctx.Form()
	user := strings.TrimSpace(form.GetString("user"))
	pwd := strings.TrimSpace(form.GetString("pwd"))
	if len(user) == 0 || len(pwd) == 0 {
		return api.ResponseWithCode(2, "缺少参数: user or pwd")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err != nil {
		return api.ResponseWithCode(3, "网络连接失败")
	}
	defer trans.Close()
	r, _ := cli.CheckLogin(thrift.Context, user, pwd, true)
	if r.ErrCode == 0 {
		memberId, _ := strconv.Atoi(r.Data["id"])
		token, _ := cli.GetToken(thrift.Context, int64(memberId), true)
		r.Data["token"] = token
		return r
	} else {
		return api.ResponseWithCode(int(r.ErrCode), r.ErrMsg)
	}
}

// 账号信息
func (m MemberApi) account(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, err1 := cli.GetAccount(thrift.Context, int64(memberId))
		if err1 == nil {
			return r
		}
		err = err1
	}
	return api.NewErrorResponse(err.Error())
}

// 账号信息
func (m MemberApi) complex(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, _ := cli.Complex(thrift.Context, memberId)
		return r
	}
	return api.NewErrorResponse(err.Error())
}

// 银行卡
func (m MemberApi) bankcard(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, _ := cli.Bankcards(thrift.Context, memberId)
		return r
	}
	return api.NewErrorResponse(err.Error())
}

// 账号信息
func (m MemberApi) profile(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, err1 := cli.GetMember(thrift.Context, memberId)
		if err1 == nil {
			return r
		}
		err = err1
	}
	return api.NewErrorResponse(err.Error())
}

// 账号信息
func (m MemberApi) checkToken(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, err1 := cli.CheckToken(thrift.Context, memberId, token)
		if err1 == nil {
			return r
		}
		err = err1
	}
	return api.NewErrorResponse(err.Error())
}

// 获取会员信息
func (m MemberApi) getMember(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		if memberId <= 0 {
			return api.NewErrorResponse("no such member")
		}
		r, _ := cli.GetMember(thrift.Context, memberId)
		return r
	}
	return api.NewErrorResponse(err.Error())
}

func (m MemberApi) receiptsCode(ctx api.Context) interface{} {
	trans, cli, _ := thrift.MemberServeClient()
	defer trans.Close()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
	arr, _ := cli.ReceiptsCodes(thrift.Context, memberId)
	mp := map[string]interface{}{
		"list":     arr,
		"provider": provider,
	}
	return mp
}

func (m MemberApi) saveReceiptsCode(ctx api.Context) interface{} {
	trans, cli, _ := thrift.MemberServeClient()
	defer trans.Close()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	data := ctx.Form().GetBytes("data")
	c := &member_service.SReceiptsCode{}
	json.Unmarshal(data, c)
	if _, ok := provider[c.Identity]; !ok {
		return api.NewErrorResponse("不支持的收款码")
	}
	memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
	r, _ := cli.SaveReceiptsCode(thrift.Context, memberId, c)
	return r
}

func (m MemberApi) toggleReceipts(ctx api.Context) interface{} {
	trans, cli, _ := thrift.MemberServeClient()
	defer trans.Close()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	identity := ctx.Form().GetString("identity")
	memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
	arr, _ := cli.ReceiptsCodes(thrift.Context, memberId)
	for _, v := range arr {
		if v.Identity == identity {
			v.State = 1 - v.State
			r, _ := cli.SaveReceiptsCode(thrift.Context, memberId, v)
			return r
		}
	}
	return api.NewErrorResponse("no such receipt code")
}

/**
 * @api {post} /member/invites 获取邀请码和邀请链接
 * @apiName invites
 * @apiGroup member
 * @apiParam {String} code 用户代码
 * @apiSuccessExample Success-Response
 * {"ErrCode":0,"ErrMsg":""9\"}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (m *MemberApi) invites(ctx api.Context) interface{} {
	trans, cli, _ := thrift.MemberServeClient()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
	member, _ := cli.GetMember(thrift.Context, memberId)
	trans.Close()
	trans2, cli2, _ := thrift.RegistryServeClient()
	defer trans2.Close()
	keys := []string{registry.Domain, registry.DomainEnabledSSL,
		registry.DomainPrefixMember,
		registry.DomainPrefixMobileMember}
	mp, _ := cli2.GetRegistries(thrift.Context, keys)
	if member != nil {
		inviteCode := member.InviteCode
		proto := types.ElseString(mp[keys[1]] == "true", "https", "http")
		// 网页推广链接
		inviteLink := fmt.Sprintf("%s://%s%s/i/%s", proto, mp[keys[2]], mp[keys[0]], inviteCode)
		// 手机网页推广链接
		mobileInviteLink := fmt.Sprintf("%s://%s%s/i/%s", proto, mp[keys[3]], mp[keys[0]], inviteCode)
		mp := map[string]string{
			"code":        inviteCode,
			"link":        inviteLink,
			"mobile_link": mobileInviteLink,
		}
		return api.NewResponse(mp)
	}
	return api.NewErrorResponse("no such user")
}

func (m MemberApi) orderSummary(ctx api.Context,memberId int64)*api.Response {
   return api.ResponseWithCode(0,"")
}


/**
 * @api {post} /member/orders_quantity 获取会员的订单状态及其数量
 * @apiGroup member
 * @apiParam {String} code 用户代码
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (m MemberApi) ordersQuantity(ctx api.Context, id int64) *api.Response {
	mp, _ := rsi.MemberService.OrdersQuantity(thrift.Context,id)
	ret := map[string]int32{
		/** 待付款订单数量 */
		"AwaitPayment":mp[int32(order.StatAwaitingPayment)],
		/** 待发货订单数量 */
		"AwaitShipment":mp[int32(order.StatAwaitingShipment)],
		/** 待收货订单数量 */
		"AwaitReceive":mp[int32(order.StatShipped)],
		/** 已完成订单数量 */
		"Completed":mp[int32(order.StatCompleted)],
	}
	return api.NewResponse(ret)
}
