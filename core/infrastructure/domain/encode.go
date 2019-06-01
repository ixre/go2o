/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-10 21:16
 * description :
 * history :
 */

package domain

import (
	"errors"
	"github.com/ixre/gof/crypto"
	"strings"
)

var (
	Sha1OffSet = ""
)

// 密码SHA1加密
func Sha1Pwd(pwd string) string {
	return crypto.Sha1([]byte(strings.Join([]string{pwd, Sha1OffSet}, "")))
}

// MD5加密
func Md5(pwd string) string {
	return crypto.Md5([]byte(pwd))
}

func ChkPwdRight(pwd string) (bool, error) {
	if len(pwd) < 6 {
		return false, errors.New("密码必须大于6位")
	}
	return true, nil
}

func Sha1(s string) string {
	return crypto.Sha1([]byte(s))
}

// 加密会员密码,因为可能会使用手机号码登录，
// 所以密码不能依据用户名作为生成凭据
func MemberSha1Pwd(pwd string) string {
	if pwd == "" {
		return ""
	}
	return Sha1Pwd(pwd)
}

// 交易密码
func TradePwd(pwd string) string {
	if pwd == "" {
		return ""
	}
	return Sha1Pwd(pwd)
}

//加密合作商密码
func MerchantSha1Pwd(user, pwd string) string {
	if pwd == "" {
		return ""
	}
	return Sha1Pwd(pwd)
}

// 密码Md5加密
func Md5Pwd(pwd, str string) string {
	return crypto.Md5([]byte(strings.Join([]string{str, Sha1OffSet, pwd}, "")))
}

// 密码SHA1加密
func ShaPwd(pwd, p string) string {
	return crypto.Sha1([]byte(strings.Join([]string{p, pwd, Sha1OffSet}, "")))
}
