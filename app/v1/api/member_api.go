package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ixre/gof/api"
	"github.com/ixre/gof/types"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/registry"
	"go2o/core/service"
	"go2o/core/service/proto"
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
	if len(code) > 0 {
		trans, cli, _ := service.MemberServiceClient()
		defer trans.Close()
		v, _ := cli.FindMember(context.TODO(),
			&proto.FindMemberRequest{
				Cred:  proto.ECredentials_Code,
				Value: code,
			})
		memberId = v.Value
	}
	switch fn {
	case "order_summary":
		return m.orderSummary(ctx, memberId)
	case "orders_quantity":
		return m.ordersQuantity(ctx, memberId)
	case "address":
		return m.address(ctx, memberId)
	case "save_address":
		return m.saveAddress(ctx, memberId)
	case "delete_address":
		return m.deleteAddress(ctx, memberId)
	case "invites":
		return m.invites(ctx, memberId)

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
	trans, cli, err := service.MemberServiceClient()
	if err != nil {
		return api.ResponseWithCode(3, "网络连接失败")
	}
	defer trans.Close()
	r, _ := cli.CheckLogin(context.TODO(), &proto.LoginRequest{
		User:   user,
		Pwd:    pwd,
		Update: true,
	})
	if r.ErrCode == 0 {
		memberId, _ := strconv.Atoi(r.Data["id"])
		token, _ := cli.GetToken(context.TODO(),
			&proto.GetTokenRequest{
				MemberId: int64(memberId),
				Reset_:   true,
			})
		r.Data["token"] = token.Value
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
	trans, cli, err := service.MemberServiceClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.FindMember(context.TODO(), &proto.FindMemberRequest{
			Cred:  proto.ECredentials_Code,
			Value: code,
		})
		r, err1 := cli.GetAccount(context.TODO(), memberId)
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
	trans, cli, err := service.MemberServiceClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.FindMember(context.TODO(),
			&proto.FindMemberRequest{
				Cred:  proto.ECredentials_Code,
				Value: code,
			})
		r, _ := cli.Complex(context.TODO(), memberId)
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
	trans, cli, err := service.MemberServiceClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.FindMember(context.TODO(),
			&proto.FindMemberRequest{
				Cred:  proto.ECredentials_Code,
				Value: code,
			})
		r, _ := cli.GetBankCards(context.TODO(), memberId)
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
	trans, cli, err := service.MemberServiceClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.FindMember(context.TODO(),
			&proto.FindMemberRequest{
				Cred:  proto.ECredentials_Code,
				Value: code,
			})
		r, err1 := cli.GetMember(context.TODO(), memberId)
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
	trans, cli, err := service.MemberServiceClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.FindMember(context.TODO(),
			&proto.FindMemberRequest{
				Cred:  proto.ECredentials_Code,
				Value: code,
			})
		r, err1 := cli.CheckToken(context.TODO(),
			&proto.CheckTokenRequest{
				MemberId: memberId.Value,
				Token:    token,
			})
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
	trans, cli, err := service.MemberServiceClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.FindMember(context.TODO(),
			&proto.FindMemberRequest{
				Cred:  proto.ECredentials_Code,
				Value: code,
			})
		if memberId.Value <= 0 {
			return api.NewErrorResponse("no such member")
		}
		r, _ := cli.GetMember(context.TODO(), memberId)
		return r
	}
	return api.NewErrorResponse(err.Error())
}

func (m MemberApi) receiptsCode(ctx api.Context) interface{} {
	trans, cli, _ := service.MemberServiceClient()
	defer trans.Close()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	memberId, _ := cli.FindMember(context.TODO(),
		&proto.FindMemberRequest{
			Cred:  proto.ECredentials_Code,
			Value: code,
		})
	arr, _ := cli.ReceiptsCodes(context.TODO(), memberId)
	mp := map[string]interface{}{
		"list":     arr,
		"provider": provider,
	}
	return mp
}

func (m MemberApi) saveReceiptsCode(ctx api.Context) interface{} {
	trans, cli, _ := service.MemberServiceClient()
	defer trans.Close()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	data := ctx.Form().GetBytes("data")
	c := &proto.SReceiptsCode{}
	_ = json.Unmarshal(data, c)
	if _, ok := provider[c.Identity]; !ok {
		return api.NewErrorResponse("不支持的收款码")
	}
	memberId, _ := cli.FindMember(context.TODO(),
		&proto.FindMemberRequest{
			Cred:  proto.ECredentials_Code,
			Value: code,
		})
	r, _ := cli.SaveReceiptsCode(context.TODO(), &proto.ReceiptsCodeSaveRequest{
		MemberId: memberId.Value,
		Code:     c,
	})
	return r
}

func (m MemberApi) toggleReceipts(ctx api.Context) interface{} {
	trans, cli, _ := service.MemberServiceClient()
	defer trans.Close()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	identity := ctx.Form().GetString("identity")
	memberId, _ := cli.FindMember(context.TODO(),
		&proto.FindMemberRequest{
			Cred:  proto.ECredentials_Code,
			Value: code,
		})
	arr, _ := cli.ReceiptsCodes(context.TODO(), memberId)
	for _, v := range arr.Value {
		if v.Identity == identity {
			v.State = 1 - v.State
			r, _ := cli.SaveReceiptsCode(context.TODO(), &proto.ReceiptsCodeSaveRequest{
				MemberId: memberId.Value,
				Code:     v,
			})
			return r
		}
	}
	return m.utils.error(errors.New("no such receipt code"))
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
func (m *MemberApi) invites(ctx api.Context, memberId int64) *api.Response {
	trans, cli, _ := service.MemberServiceClient()
	member, _ := cli.GetMember(context.TODO(), &proto.Int64{Value: memberId})
	trans.Close()
	trans2, cli2, _ := service.RegistryServiceClient()
	defer trans2.Close()
	mp1, _ := cli2.GetValue(context.TODO(), &proto.String{Value: registry.Domain})
	mp2, _ := cli2.GetValue(context.TODO(), &proto.String{Value: registry.HttpProtocols})
	mp3, _ := cli2.GetValue(context.TODO(), &proto.String{Value: registry.DomainPrefixMember})
	mp4, _ := cli2.GetValue(context.TODO(), &proto.String{Value: registry.DomainPrefixMobileMember})
	trans.Close()
	domain := mp1.Value
	proto := types.ElseString(mp2.Value == "true", "https", "http")

	if member != nil {
		inviteCode := member.InviteCode
		// 网页推广链接
		inviteLink := fmt.Sprintf("%s://%s%s/i/%s",
			proto,
			mp3.Value,
			domain,
			inviteCode)
		// 手机网页推广链接
		mobileInviteLink := fmt.Sprintf("%s://%s%s/i/%s",
			proto,
			mp4.Value,
			domain,
			inviteCode)
		mp := map[string]string{
			"code":        inviteCode,
			"link":        inviteLink,
			"mobile_link": mobileInviteLink,
		}
		return m.utils.success(mp)
	}
	return m.utils.error(errors.New("no such user"))
}

func (m MemberApi) orderSummary(ctx api.Context, memberId int64) *api.Response {
	return api.ResponseWithCode(0, "")
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
	trans, cli, _ := service.MemberServiceClient()
	defer trans.Close()
	mp, _ := cli.OrdersQuantity(context.TODO(), &proto.Int64{Value: id})
	ret := map[string]int32{
		/** 待付款订单数量 */
		"AwaitPayment": mp.Data[int32(order.StatAwaitingPayment)],
		/** 待发货订单数量 */
		"AwaitShipment": mp.Data[int32(order.StatAwaitingShipment)],
		/** 待收货订单数量 */
		"AwaitReceive": mp.Data[int32(order.StatShipped)],
		/** 已完成订单数量 */
		"Completed": mp.Data[int32(order.StatCompleted)],
	}
	return m.utils.success(ret)
}

/**
 * @api {post} /member/address 获取会员的收货地址
 * @apiGroup member
 * @apiParam {String} code 用户代码
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (m MemberApi) address(ctx api.Context, memberId int64) *api.Response {
	trans, cli, _ := service.MemberServiceClient()
	defer trans.Close()
	address, _ := cli.GetAddressList(context.TODO(), &proto.Int64{Value: memberId})
	return m.utils.success(address)
}

/**
 * @api {post} /member/save_address 保存会员的收货地址
 * @apiGroup member
 * @apiParam {String} code 用户代码
 * @apiParam {int} address_id 地址编号, 保存时需传递
 * @apiParam {int} consignee_name 收货人姓名
 * @apiParam {int} consignee_phone 收货人电话
 * @apiParam {int} province 数字编码(省)
 * @apiParam {int} city 数字编码(市)
 * @apiParam {int} district 数字编码(区)
 * @apiParam {int} detail_address 详细的地址,比如: 幸福路12号
 * @apiParam {int} is_default 是否默认收货地址, 1:表示默认收货地址
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (m MemberApi) saveAddress(ctx api.Context, memberId int64) *api.Response {
	form := ctx.Form()
	var e = proto.SAddress{
		ID:             int64(form.GetInt("address_id")),
		ConsigneeName:  form.GetString("consignee_name"),
		ConsigneePhone: form.GetString("consignee_phone"),
		Province:       int32(form.GetInt("province")),
		City:           int32(form.GetInt("city")),
		District:       int32(form.GetInt("district")),
		DetailAddress:  form.GetString("detail_address"),
		IsDefault:      int32(form.GetInt("is_default")),
	}
	trans, cli, _ := service.MemberServiceClient()
	defer trans.Close()
	ret, _ := cli.SaveAddress(context.TODO(), &proto.SaveAddressRequest{
		MemberId: memberId,
		Value:    &e,
	})
	return api.NewResponse(ret)
}

/**
 * @api {post} /member/delete_address 删除会员的收货地址
 * @apiGroup member
 * @apiParam {String} code 用户代码
 * @apiParam {int} address_id 地址编号
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (m MemberApi) deleteAddress(ctx api.Context, memberId int64) *api.Response {
	addressId := int64(ctx.Form().GetInt("address_id"))
	trans, cli, _ := service.MemberServiceClient()
	defer trans.Close()
	ret, _ := cli.DeleteAddress(context.TODO(), &proto.AddressIdRequest{
		MemberId:  memberId,
		AddressId: addressId,
	})
	return m.result(ret)
}
