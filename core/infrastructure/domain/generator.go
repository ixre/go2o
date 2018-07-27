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

// 创建交易号(16位)，business为零时，交易号为15位
func NewTradeNo(business int, userId int) string {
	dt := time.Now()
	rand.Seed(dt.UnixNano())
	second := 1000 + (dt.Minute()*60)*(dt.Hour()/12+1) + dt.Second()
	arr := make([]string, 5)
	if business > 0 && business < 10 {
		arr[0] = strconv.Itoa(business) // 业务：长度1
	}
	arr[1] = dt.Format("060102")              // 年月日：长度6
	arr[2] = strconv.Itoa(userId)             // 用户编号:后3位
	arr[3] = strconv.Itoa(second)             // 秒:4位
	arr[4] = strconv.Itoa(10 + rand.Intn(88)) // 随机数:2位
	// 将用户编号调整为3位
	if l := len(arr[2]); l > 3 {
		arr[2] = arr[2][l-3:]
	} else if l < 3 {
		arr[2] = strings.Repeat("0", 3-l) + arr[2]
	}
	return strings.Join(arr, "")
}

// 创建邀请码(6位)
func GenerateInvitationCode() string {
	var seed = fmt.Sprintf("%d%s", time.Now().Unix(), util.RandString(6))
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

// 创建随机数字密码
func GenerateRandomIntPwd(n int) string {
	return strconv.Itoa(util.RandInt(n))
}

// 生成16位唯一的md5购物车码
func GenerateCartCode(unix int64, nano int) string {
	str := fmt.Sprintf("%d-%d*%d-%d", unix, nano,
		unix%int64(nano), util.RandInt(3))
	result := crypto.Md5([]byte(str))
	return result[8:24]
}
