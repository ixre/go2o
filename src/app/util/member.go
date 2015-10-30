/**
 * Copyright 2015 @ z3q.net.
 * name : member.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package util

import (
	"errors"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/crypto"
	"github.com/jsix/gof/web"
	"strconv"
)

const offset string = "%$^&@#"

func chkStorage(sto gof.Storage) {
	if sto == nil {
		panic(errors.New("[ Api] - api token storage is null !"))
	}
}

// 获取会员API调用密钥Key
func GetMemberApiTokenKey(memberId int) string {
	return fmt.Sprintf("api:member:token:%d", memberId)
}

// 设置令牌，并返回
func SetMemberApiToken(sto gof.Storage, memberId int, pwd string) string {
	chkStorage(sto)
	cyp := crypto.NewUnixCrypto(pwd+offset, offset)
	var token string = string(cyp.Encode())
	var key string = GetMemberApiTokenKey(memberId)

	sto.Set(key, token)      // 存储令牌
	sto.Set(key+"base", pwd) // 存储令牌凭据

	return token
}

// 获取会员的API令牌
func GetMemberApiToken(sto gof.Storage, memberId int) (string, string) {
	chkStorage(sto)

	var key = GetMemberApiTokenKey(memberId)
	var srcToken, tokenBase string

	srcToken, _ = sto.GetString(key)
	tokenBase, _ = sto.GetString(key + "base")
	return srcToken, tokenBase
}

// 移除会员令牌
func RemoveMemberApiToken(sto gof.Storage, memberId int, token string) bool {
	srcToken, _ := GetMemberApiToken(sto, memberId)
	if len(srcToken) == 0 && srcToken == token {
		var key string = GetMemberApiTokenKey(memberId)
		sto.Del(key)
		sto.Del(key + "base")

	}
	return false
}

// 校验令牌
func CompareMemberApiToken(sto gof.Storage, memberId int, token string) bool {

	if len(token) == 0 {
		return false
	}

	srcToken, tokenBase := GetMemberApiToken(sto, memberId)
	if len(srcToken) == 0 || len(tokenBase) == 0 {
		return false
	}

	//fmt.Println("-----",srcToken," IN:",token , " equals:",srcToken==token, len(srcToken),len(token))
	return srcToken == token
}

// 会员Http请求会话链接
func MemberHttpSessionConnect(ctx *web.Context, call func(memberId int)) (ok bool, memberId int) {
	//return true,30
	// 如果传递会话参数正确，能存储到Session

	form := ctx.Request.URL.Query()
	if memberId, err := strconv.Atoi(form.Get("member_id")); err == nil {
		var token string = form.Get("token")
		if CompareMemberApiToken(ctx.App.Storage(), memberId, token) {
			if call != nil {
				call(memberId)
			}
			ctx.Session().Set("client_member_id", memberId)
			ctx.Session().Save()
			return true, memberId
		}
	} else {
		// 如果没有传递参数从会话中获取
		if v := ctx.Session().Get("client_member_id"); v != nil {
			memberId = v.(int)
			return true, memberId
		}
	}

	//http://zs.ts.com/main/msc?device=1&return_url=/list/all_cate&member_id=30&token=25245e2640232df15db617473f59159c9d3d7c300ce349cb9a953b
	//SetMemberApiToken(ctx.App.Storage(),30,"f22e180335baf50c134ea5c1093de0a6")
	//fmt.Println(GetMemberApiToken(ctx.App.Storage(),30))

	return false, memberId
}

// 会员Http请求会话链接
func MemberHttpSessionDisconnect(ctx *web.Context) bool {
	form := ctx.Request.URL.Query()
	if memberId, err := strconv.Atoi(form.Get("member_id")); err == nil {
		var token string = form.Get("token")
		return RemoveMemberApiToken(ctx.App.Storage(), memberId, token)
	}
	return false
}
