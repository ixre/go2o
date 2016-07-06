/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package domain

import (
	"fmt"
	"github.com/jsix/gof/crypto"
	"github.com/jsix/gof/util"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//新订单号
func NewOrderNo(vendorId int) string {
	//MerchantId的首位和末尾再加7位随机数
	unix := time.Now().UnixNano()
	rand.Seed(unix)
	typeStr := ""
	timeStr := time.Now().Format("0601")
	rd := strconv.Itoa(time.Now().Nanosecond() + rand.Intn(999-100))
	if l := len(rd); l > 6 {
		rd = rd[:6]
	} else {
		rd = strings.Repeat("0", 6-l) + rd
	}
	if vendorId > 0 {
		typeStr = "MC"
	} else {
		typeStr = "ZY"
	}
	return fmt.Sprintf("%s%s%s", typeStr, timeStr, rd)
}

// 新交易号(12位)
func NewTradeNo(merchantId int) string {
	unix := time.Now().UnixNano()
	rand.Seed(unix)
	rd := 10000 + rand.Intn(9999-1000)
	timeStr := time.Now().Format("0602")
	ptStr := strconv.Itoa(merchantId)
	return fmt.Sprintf("%s%s%d", ptStr, timeStr, rd)
}

// 创建邀请码(6位)
func GenerateInvitationCode() string {
	var seed string = fmt.Sprintf("%d%s", time.Now().Unix(), util.RandString(6))
	var md5 = crypto.Md5([]byte(seed))
	return md5[8:16]
}

// 获取新的验证码
func NewCheckCode() string {
	unix := time.Now().UnixNano()
	rand.Seed(unix)
	rd := 1000 + rand.Intn(9999-1000)
	return strconv.Itoa(rd)
}

// 创建API编号(10位)
func NewApiId(id int) string {
	var offset = id*360 + id%2
	return fmt.Sprintf("60%s%d", strings.Repeat("0", 8-len(strconv.Itoa(offset))), offset)
}

//创建密钥(16位)
func NewSecret(hex int) string {
	str := fmt.Sprintf("%d$%d", hex, time.Now().Add(time.Hour*24*365).Unix())
	return crypto.Md5([]byte(str))[8:24]
}

// 创建随机密码
func GenerateRandomPwd(n int) string {
	return util.RandString(n)
}

// 创建随机数字密码
func GenerateRandomIntPwd(n int) string {
	return strconv.Itoa(util.RandInt(n))
}
