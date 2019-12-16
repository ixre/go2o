package api

import (
	"github.com/ixre/gof/api"
	"go2o/core/domain/interface/member"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"strings"
)

var _ api.Handler = new(accountApi)

type accountApi struct {
	utils
}

func (a accountApi) Process(fn string, ctx api.Context) *api.Response {
	var memberId int64
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) > 0 {
		memberId, _ = rsi.MemberService.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
	}
	switch fn {
	case "wallet_log":
		return a.WalletAccountLog(ctx, memberId)
	case "integral_log":
		return a.integralAccountLog(ctx, memberId)
	case "balance_log":
		return a.balanceAccountLog(ctx, memberId)
	}
	return api.ResponseWithCode(-1,"api not defined")
}

func (a accountApi) accountLog(ctx api.Context, memberId int64,account int) *api.Response {
	begin := int32(ctx.Form().GetInt("begin"))
	size := int32(ctx.Form().GetInt("size"))
	p := &ttype.SPagingParams{
		SortBy:"create_time DESC,bi.id DESC",
		Begin:      begin,
		Over:       begin+size,
	}
	r,_ := rsi.MemberService.PagingAccountLog(thrift.Context,int64(memberId), int32(account),p)
	return api.NewResponse(r)
}

/**
 * @api {post} /account/wallet_log 获取会员的钱包明细
 * @apiGroup account
 * @apiParam {String} code 用户代码
 * @apiParam {Int} begin 开始记录数,默认为:0
 * @apiParam {Int} size 数量
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (a accountApi) WalletAccountLog(ctx api.Context, memberId int64) *api.Response {
	return a.accountLog(ctx,memberId,member.AccountWallet)
}

/**
 * @api {post} /account/wallet_log 获取会员的钱包明细
 * @apiGroup account
 * @apiParam {String} code 用户代码
 * @apiParam {Int} begin 开始记录数,默认为:0
 * @apiParam {Int} size 数量
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (a accountApi) walletAccountLog(ctx api.Context, memberId int64) *api.Response {
	return a.accountLog(ctx,memberId,member.AccountWallet)
}

/**
 * @api {post} /account/integral_log 获取会员的积分明细
 * @apiGroup account
 * @apiParam {String} code 用户代码
 * @apiParam {Int} begin 开始记录数,默认为:0
 * @apiParam {Int} size 数量
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (a accountApi) integralAccountLog(ctx api.Context, memberId int64) *api.Response {
	return a.accountLog(ctx, memberId,member.AccountIntegral)

}

/**
 * @api {post} /account/balance_log 获取会员的余额明细
 * @apiGroup account
 * @apiParam {String} code 用户代码
 * @apiParam {Int} begin 开始记录数,默认为:0
 * @apiParam {Int} size 数量
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (a accountApi) balanceAccountLog(ctx api.Context, memberId int64) *api.Response {
	return a.accountLog(ctx, memberId,member.AccountBalance)
}