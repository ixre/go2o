/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-10 21:16
 * description :
 * history :
 */

package domain

import (
	"fmt"
	"ops/cf/crypto"
	"strings"
	"time"
)

//加密会员密码
func EncodeMemberPwd(usr, pwd string) string {
	return crypto.Md5([]byte(strings.Join([]string{usr, "$OPSoft$", pwd}, "")))
}

//加密合作商密码
func EncodePartnerPwd(usr, pwd string) string {
	return crypto.Md5([]byte(strings.Join([]string{usr, "$OPSoft$", pwd}, "")))
}

//创建密钥
func NewSecret(hex int) string {
	str := fmt.Sprintf("%d$%d", hex, time.Now().Add(time.Hour*24*365).Unix())
	return crypto.Md5([]byte(str))[8:24]
}
