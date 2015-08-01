/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-10 21:16
 * description :
 * history :
 */

package domain

import (
	"errors"
	"github.com/atnet/gof/crypto"
	"strings"
)

func ChkPwdRight(pwd string) (bool, error) {
	if len(pwd) < 6 {
		return false, errors.New("密码必须大于6位")
	}
	return true, nil
}

// 加密会员密码,因为可能会使用手机号码登陆，
// 所以密码不能依据用户名作为生成凭据
func MemberSha1Pwd(pwd string) string {
	return crypto.Sha1([]byte(ShaPwd(pwd, "")))
}

//加密合作商密码
func Md5PartnerPwd(usr, pwd string) string {
	return Md5Pwd(pwd, usr)
}

// 密码Md5加密
func Md5Pwd(pwd, offset string) string {
	return crypto.Md5([]byte(strings.Join([]string{offset, "go2o@S1N1.COM", pwd}, "")))
}

// 密码SHA1加密
func ShaPwd(pwd, offset string) string {
	return crypto.Sha1([]byte(strings.Join([]string{offset, "h3f", pwd, ".net"}, "")))
}
