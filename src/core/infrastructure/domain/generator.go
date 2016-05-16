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

const (
	minRand int = 100000
	maxRand int = 999900
)

//新订单号
func NewOrderNo(partnerId int) string {
	//MerchantId的首位和末尾再加7位随机数
	unix := time.Now().UnixNano()
	rand.Seed(unix)
	rd := minRand + rand.Intn(maxRand-minRand) //minRand - maxRand中间的随机数
	timeStr := time.Now().Format("0601")
	ptStr := strconv.Itoa(partnerId)
	return fmt.Sprintf("%s%s%s%d", ptStr[:1], timeStr, ptStr[len(ptStr)-1:], rd)
}

// 新交易号(12位)
func NewTradeNo(partnerId int) string {
	unix := time.Now().UnixNano()
	rand.Seed(unix)
	rd := 10000 + rand.Intn(9999-1000)
	timeStr := time.Now().Format("0602")
	ptStr := strconv.Itoa(partnerId)
	return fmt.Sprintf("%s%s%d", ptStr, timeStr, rd)
}

// 创建邀请码(6位)
func GenerateInvitationCode() string {
	var seed string = fmt.Sprintf("%d%s", time.Now().Unix(), util.RandString(6))
	var md5 = crypto.Md5([]byte(seed))
	return md5[8:14]
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
