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
	"github.com/atnet/gof/crypto"
	"strings"
)

func ChkPwdRight(pwd string) (bool, error) {
	return true, nil
}

//加密会员密码
func Md5MemberPwd(usr, pwd string) string {
	return Md5Pwd(pwd,"member_"+usr)
}

//加密合作商密码
func Md5PartnerPwd(usr, pwd string) string {
	return Md5Pwd(pwd, usr)
}

// 密码Md5加密
func Md5Pwd(pwd, offset string) string {
	return crypto.Md5([]byte(strings.Join([]string{offset, "go2o@S1N1.COM", pwd}, "")))
}
