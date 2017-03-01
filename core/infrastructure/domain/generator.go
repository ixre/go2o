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
	"bytes"
	"fmt"
	"github.com/jsix/gof/crypto"
	"github.com/jsix/gof/storage"
	"github.com/jsix/gof/util"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//新订单号
func NewOrderNo(vendorId int, prefix string) string {
	rdLen := 6 - len(prefix)
	//MerchantId的首位和末尾再加7位随机数
	unix := time.Now().UnixNano()
	rand.Seed(unix)
	buf := bytes.NewBufferString(prefix)
	vendorStr := strconv.Itoa(vendorId)
	if l := len(vendorStr); l < 6 {
		buf.WriteString("1")
		buf.WriteString(strings.Repeat("0", 5-l))
	}
	buf.WriteString(vendorStr)
	rd := strconv.Itoa(rand.Intn(999999 - 100000))
	if l := len(rd); l > rdLen {
		buf.WriteString(rd[:rdLen])
	} else {
		buf.WriteString(strings.Repeat("0", rdLen-l))
		buf.WriteString(rd)
	}
	return buf.String()
}

// 新交易号(12位)
func NewTradeNo(mchId int) string {
	unix := time.Now().UnixNano()
	rand.Seed(unix)
	rd := 10000 + rand.Intn(9999-1000)
	timeStr := time.Now().Format("0602")
	ptStr := strconv.Itoa(mchId)
	return fmt.Sprintf("%s%s%d", ptStr, timeStr, rd)
}

// 创建邀请码(6位)
func GenerateInvitationCode() string {
	var seed string = fmt.Sprintf("%d%s", time.Now().Unix(), util.RandString(6))
	var md5 = crypto.Md5([]byte(seed))
	return md5[8:14]
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

// 生成不重复的交易号(通过存储)
func NewTradeNoFromStorage(s storage.Interface, prefix string) string {
	for {
		no := crypto.Md5([]byte(NewTradeNo(0) + GenerateRandomIntPwd(10)))
		if !s.Exists(prefix + no) {
			return no
		}
	}
	return ""
}

// 生成16位唯一的md5购物车码
func GenerateCartCode(unix int64, nano int) string {
	str := fmt.Sprintf("%d-%d*%d-%d", unix, nano,
		unix%int64(nano), util.RandInt(1000))
	result := crypto.Md5([]byte(str))
	return result[8:24]
}
